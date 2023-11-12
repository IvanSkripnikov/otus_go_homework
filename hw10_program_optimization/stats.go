package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	var user User
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)

	for i := 0; scanner.Scan(); i++ {
		if err := user.UnmarshalJSON(scanner.Bytes()); err != nil {
			return nil, fmt.Errorf("read user error: %w", err)
		}

		if !strings.HasSuffix(user.Email, domain) {
			continue
		}
		emails := strings.SplitN(user.Email, "@", 2)

		if len(emails) > 1 {
			domainKey := strings.ToLower(emails[1])
			result[domainKey]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read users error: %w", err)
	}

	return result, nil
}
