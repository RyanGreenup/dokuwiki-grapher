package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type node struct {
	id    string
	label string
}

type edge struct {
	from string
	to   string
}

func main() {
	HOME, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Unable to get user home directory")
	}
	dokuwiki_dir := HOME + "/Notes/dokuwiki/data/pages" // DIRECTORY HERE
	os.Chdir(dokuwiki_dir)
	files := get_txt_files(".")

	nodes := []node{}
	edges := []edge{}

	nodes = add_nodes(nodes, files)
	edges = get_links(edges, files)


    fmt.Println("The Nodes are")
	fmt.Println(nodes)

    fmt.Println("The Edges are")
	fmt.Println(edges)

}

func add_nodes(nodes []node, files []string) []node {
	for _, file := range files {
		file = filepath.Clean(file)
		// TODO this could be more elegant with a reverse and replaceAll
		file = strings.Replace(file, `.txt`, ``, 1)
		file = strings.ReplaceAll(file, `/`, `:`)
		nodes = append(nodes, node{id: file, label: file})
	}

	return nodes
}

func get_links(edges []edge, files []string) []edge {
	re, err := regexp.Compile(`\[\[(.*)\]\]`)
	if err != nil {
		log.Fatal("E: 2891, Unable to compile Regex")
	}
	for i, file_path := range files {
		_ = i

		file, err := os.Open(file_path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
        buf := make([]byte, 0, 64*1024)
        scanner.Buffer(buf, 1024*1024*1024) // TODO BAD HACK, why is there a 1GB line?
		// optionally resize scanner's capacity for lines over 64k
		for scanner.Scan() {
			text := scanner.Text()
			matches := re.FindAllStringSubmatch(text, -1)
			if len(matches) > 0 {

				pagename := strings.Replace(file_path, `.txt`, ``, 1)
				pagename = strings.ReplaceAll(pagename, `/`, `:`)
				m := matches[0][1]
				m = clean_link(m)
				edges = append(edges, edge{from: pagename, to: m})
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

	}
	return edges
}

func clean_link(link_text string) string {
	// Try with space first
	re := regexp.MustCompile(` \|.*`)
	link_text = re.ReplaceAllString(link_text, ``) // Remove everything after pipe
	// Now without space
	re = regexp.MustCompile(`\|.*`)
	link_text = re.ReplaceAllString(link_text, ``) // Remove everything after pipe

	link_text = strings.ReplaceAll(link_text, ` `, `_`)
	link_text = strings.ToLower(link_text)
	return link_text
}

func get_txt_files(top_dir string) []string {
	ext := ".txt" // EXTENSION HERE

	files := []string{}
	err := filepath.Walk(top_dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(path) == ext {
				path, err = filepath.Rel(top_dir, path)
				if err != nil {
					log.Println("Unable to get relative path")
				}
				files = append(files, path)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return files
}
