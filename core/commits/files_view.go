package commits

import "fmt"

func GitStashView() string {
	// get git status
	status, gitChangeStatus, err := getGitStatus()
	if err != nil {
		return ""
	}
	fmt.Println(status)
	fmt.Println(gitChangeStatus)
	return ""
}
