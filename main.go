package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/google/go-github/v55/github"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(fmt.Sprintf("Usage: %s <github api token>", os.Args[0]))
		os.Exit(1)
	}

	for i := 1; i < len(os.Args); i++ {
		client := github.NewClient(nil).WithAuthToken(os.Args[i])

		re := regexp.MustCompile("https://github.com/([^/]+)/([^/]+)/pull/(\\d+)")

		// read ~/chronomerge.txt, each line contains a pull request url
		// for each pull request, check if it can be merged and if it can
		// merge it

		f, err := os.Open(os.Getenv("HOME") + "/chronomerge.txt")
		if err != nil {
			panic(err)
		}

		output := []string{}

		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			// parse the url with regex https://github.com/([^/]+)/([^/]+)/pull/(\d+)
			// get the owner, repo and pull number

			line := scanner.Text()

			match := re.FindStringSubmatch(line)

			if len(match) != 4 {
				fmt.Println("Invalid url:", scanner.Text())
				continue
			}

			owner := match[1]
			repo := match[2]
			pullNumber, err := strconv.ParseInt(match[3], 10, 64)
			if err != nil {
				fmt.Println("Invalid url:", scanner.Text())
				continue
			}

			// merge it
			for _, strategy := range []string{"squash", "merge"} {
				_, _, err = client.PullRequests.Merge(context.Background(), owner, repo, int(pullNumber), "", &github.PullRequestOptions{
					MergeMethod: strategy,
				})
				if err != nil {
					fmt.Println("Error merging pull request:", err)
					// retain the url
					output = append(output, line)
					continue
				}
				break
			}
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}

		if err := f.Close(); err != nil {
			panic(err)
		}

		// write the remaining urls to ~/chronomerge.txt
		f, err = os.Create(os.Getenv("HOME") + "/chronomerge.txt")
		if err != nil {
			panic(err)
		}

		for _, line := range output {
			fmt.Fprintln(f, line)
		}

		if err := f.Close(); err != nil {
			panic(err)
		}
	}
}
