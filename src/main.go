package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"runtime"

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
		cmd = exec.Command("open", url) // "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = exec.Command("xdg-open", url) //  "xdg-open"
	}

	err := cmd.Run()

	if err != nil {
		log_error(err, 1)
	}
}

func handleOpen(cmd *cobra.Command, args []string) {

	var domains []Domain
	var err error

	if len(args) == 0 {
		domains, err = queryDomains("")
	} else {
		domains, err = queryDomains(args[0])
	}
	if err != nil {
		log_error(err, 1)
	}

	var suggestions []string

	// Loop through the domains and append the Value of each to the values slice
	for _, domain := range domains {
		suggestions = append(suggestions, fmt.Sprintf("%v", domain.Value))
	}

	if len(suggestions) == 0 {
		log_error(errors.New("cannot open an empty index, call `add` first"), 0)
		return
	}

	// The question to ask
	var qs = []*survey.Question{
		{
			Name: "item",
			Prompt: &survey.Select{
				Message: "Choose a bookmark to Open:",
				Options: suggestions,
			},
		},
	}

	// The answer will be stored in this struct
	answer := struct {
		Item string `survey:"item"` // matches the question name
	}{}

	// Perform the survey
	err = survey.Ask(qs, &answer)
	if err != nil {
		log_error(err, 1)
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

	err := addDomain(args[0], aliasFlag)
	if err != nil {
		log_error(err, 1)
	} else {
		fmt.Printf("Successfully bookmarked %v!", args[0])
	}
}

func handleDelete(cmd *cobra.Command, args []string) {

	var domains []Domain
	var err error

	if len(args) == 0 {
		domains, err = queryDomains("")
	} else {
		domains, err = queryDomains(args[0])
	}
	if err != nil {
		log_error(err, 1)
	}

	var suggestions []string

	// Loop through the domains and append the Value of each to the values slice
	for _, domain := range domains {
		suggestions = append(suggestions, domain.Value)
	}

	if len(suggestions) == 0 {
		fmt.Println("Cannot delete in an empty index")
		return
	}

	// The question to ask
	var qs = []*survey.Question{
		{
			Name: "item",
			Prompt: &survey.Select{
				Message: "Choose a bookmark to delete:",
				Options: suggestions,
			},
		},
	}

	// The answer will be stored in this struct
	answer := struct {
		Item string `survey:"item"` // matches the question name
	}{}

	// Perform the survey
	err = survey.Ask(qs, &answer)
	if err != nil {
		log_error(err, 1)
		return
	}

	// resolve delete
	var id string
	for _, domain := range domains {
		if domain.Value == answer.Item {
			id = domain.ID
		}
	}
	err = deleteDomain(id)
	if err != nil {
		log_error(err, 1)
	} else {
		log.Printf("Successfully deleted %v", answer.Item)

	}
}

func handleImport(cmd *cobra.Command, args []string) {
	err := importer(args[0])
	if err != nil {
		return
	}
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
		Use:   "open [path]",
		Short: "Open a bookmark",
		Args:  cobra.MinimumNArgs(0),
		Run:   handleOpen,
	}

	var cmdDelete = &cobra.Command{
		Use:   "delete [path]",
		Short: "Delete a bookmark",
		Args:  cobra.MinimumNArgs(0),
		Run:   handleDelete,
	}

	var cmdImport = &cobra.Command{
		Use:   "import [path]",
		Short: "Import your bookmarks for an exported HTML file",
		Args:  cobra.MinimumNArgs(1),
		Run:   handleImport,
	}
	rootCmd.AddCommand(cmdAdd, cmdOpen, cmdDelete, cmdImport)

	if err := rootCmd.Execute(); err != nil {
		log_error(err, 2)
	}
}
