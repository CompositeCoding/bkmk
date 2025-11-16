package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"runtime"

	"bkmk/importer"
	"bkmk/lib"
	"bkmk/store"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

func isValidURL(url string) bool {
	regex := regexp.MustCompile(`^((http|https|ftp|file):\/\/)?[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?$`)
	return regex.MatchString(url)
}

func openBrowser(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll", "FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	if err := cmd.Run(); err != nil {
		lib.LogError(err, 1)
	}
}

func handleOpen(cmd *cobra.Command, args []string) {
	var domains []store.Domain
	var err error

	if len(args) == 0 {
		domains, err = store.QueryDomains("")
	} else {
		domains, err = store.QueryDomains(args[0])
	}
	if err != nil {
		lib.LogError(err, 1)
		return
	}

	var suggestions []string
	for _, d := range domains {
		suggestions = append(suggestions, fmt.Sprintf("%v", d.Value))
	}

	if len(suggestions) == 0 {
		lib.LogError(errors.New("cannot open an empty index, call `add` first"), 0)
		return
	}

	var qs = []*survey.Question{
		{
			Name: "item",
			Prompt: &survey.Select{
				Message: "Choose a bookmark to Open:",
				Options: suggestions,
			},
		},
	}

	answer := struct {
		Item string `survey:"item"`
	}{}

	if err := survey.Ask(qs, &answer); err != nil {
		lib.LogError(err, 1)
		return
	}

	openBrowser(answer.Item)
}

func handleAdd(cmd *cobra.Command, args []string) {
	if !isValidURL(args[0]) {
		fmt.Printf("Error! Unable to add %v, bkmk only supports valid URLs", args[0])
		return
	}

	aliasFlag, _ := cmd.Flags().GetString("alias")

	if err := store.AddDomain(args[0], aliasFlag); err != nil {
		lib.LogError(err, 1)
	} else {
		fmt.Printf("Successfully bookmarked %v!", args[0])
	}
}

func handleDelete(cmd *cobra.Command, args []string) {
	var domains []store.Domain
	var err error

	if len(args) == 0 {
		domains, err = store.QueryDomains("")
	} else {
		domains, err = store.QueryDomains(args[0])
	}
	if err != nil {
		lib.LogError(err, 1)
		return
	}

	var suggestions []string
	for _, d := range domains {
		suggestions = append(suggestions, d.Value)
	}

	if len(suggestions) == 0 {
		fmt.Println("Cannot delete in an empty index")
		return
	}

	var qs = []*survey.Question{
		{
			Name: "item",
			Prompt: &survey.Select{
				Message: "Choose a bookmark to delete:",
				Options: suggestions,
			},
		},
	}

	answer := struct {
		Item string `survey:"item"`
	}{}

	if err := survey.Ask(qs, &answer); err != nil {
		lib.LogError(err, 1)
		return
	}

	var id string
	for _, d := range domains {
		if d.Value == answer.Item {
			id = d.ID
			break
		}
	}

	if id == "" {
		lib.LogError(errors.New("could not resolve selected bookmark id"), 1)
		return
	}

	if err := store.DeleteDomain(id); err != nil {
		lib.LogError(err, 1)
	} else {
		log.Printf("Successfully deleted %v", answer.Item)
	}
}

func handleImport(cmd *cobra.Command, args []string) {
	if err := importer.Importer(args[0]); err != nil {
		lib.LogError(err, 1)
		return
	}
}

func handleProfile(cmd *cobra.Command, args []string) {
	// TODO: implement profiles
}

func main() {
	var rootCmd = &cobra.Command{Use: "bkmk"}

	var cmdAdd = &cobra.Command{
		Use:   "add [string to add]",
		Short: "Add a new url to your bookmarks",
		Args:  cobra.MinimumNArgs(1),
		Run:   handleAdd,
	}
	cmdAdd.Flags().StringP("alias", "a", "", "Set an Alias for the bookmark")

	var cmdOpen = &cobra.Command{
		Use:   "open [url]",
		Short: "Open a bookmark",
		Args:  cobra.MinimumNArgs(0),
		Run:   handleOpen,
	}

	var cmdDelete = &cobra.Command{
		Use:   "delete [url]",
		Short: "Delete a bookmark",
		Args:  cobra.MinimumNArgs(0),
		Run:   handleDelete,
	}

	var cmdImport = &cobra.Command{
		Use:   "import [path]",
		Short: "Import your bookmarks from a browser exported HTML file",
		Args:  cobra.MinimumNArgs(1),
		Run:   handleImport,
	}

	var cmdProfile = &cobra.Command{
		Use:   "profile [name]",
		Short: "Create or switch to a bookmark profile",
		Args:  cobra.MinimumNArgs(1),
		Run:   handleProfile,
	}
	cmdProfile.Flags().StringP("create", "c", "", "Create a new profile")
	cmdProfile.Flags().StringP("switch", "s", "", "Switch to a profile")

	rootCmd.AddCommand(cmdAdd, cmdOpen, cmdDelete, cmdImport, cmdProfile)

	if err := rootCmd.Execute(); err != nil {
		lib.LogError(err, 2)
	}
}
