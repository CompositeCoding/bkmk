package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

func isValidURL(url string) bool {
	regex := regexp.MustCompile(`^(http|https):\/\/[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?$`)
	return regex.MatchString(url)
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
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
		log.Fatal("Fatal query error")
	}

	var suggestions []string

	// Loop through the domains and append the Value of each to the values slice
	for _, domain := range domains {
		suggestions = append(suggestions, domain.Value)
	}

	if len(suggestions) == 0 {
		fmt.Println("Cannot open an empty index, call `add` first!")
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
		fmt.Println(err.Error())
		return
	}

	openBrowser(answer.Item)

}

func handleAdd(cmd *cobra.Command, args []string) {
	if !isValidURL(args[0]) {
		fmt.Printf("Error! Unable to add %v, bkmk only supports valid URLs", args[0])
		return
	}

	err := addDomain(args[0])
	if err != nil {
		fmt.Printf("Error! Unable to add %v", args[0])
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
		log.Fatal("Fatal query error")
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
				Message: "Choose a bookmark to open:",
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
		fmt.Println(err.Error())
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
		log.Print(err)
	} else {
		log.Printf("Successfully deleted %v", answer.Item)
	}
}

func handleImport(cmd *cobra.Command, args []string) {
	importer(args[0])
}

func main() {

	var rootCmd = &cobra.Command{Use: "bkmk"}

	var cmdAdd = &cobra.Command{
		Use:   "add [string to add]",
		Short: "Add a new url to your bookmarks",
		Args:  cobra.MinimumNArgs(1),
		Run:   handleAdd,
	}

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
		log.Fatal(err)
		os.Exit(1)
	}
}
