// Package portgroup groups ports by protocol and service category.
package portgroup

import "fmt"

// Category represents a logical grouping of ports.
type Category string

const (
	CategoryWeb      Category = "web"
	CategoryDatabase Category = "database"
	CategoryRemote   Category = "remote"
	CategoryMail     Category = "mail"
	CategoryDNS      Category = "dns"
	CategoryOther    Category = "other"
)

// Port represents a port with protocol.
type Port struct {
	Number   int
	Protocol string
}

// Grouper assigns ports to categories.
type Grouper struct {
	custom map[int]Category
}

var defaultCategories = map[int]Category{
	80:   CategoryWeb,
	443:  CategoryWeb,
	8080: CategoryWeb,
	8443: CategoryWeb,
	3306: CategoryDatabase,
	5432: CategoryDatabase,
	6379: CategoryDatabase,
	27017: CategoryDatabase,
	22:   CategoryRemote,
	23:   CategoryRemote,
	3389: CategoryRemote,
	25:   CategoryMail,
	465:  CategoryMail,
	587:  CategoryMail,
	993:  CategoryMail,
	53:   CategoryDNS,
}

// New returns a Grouper with optional custom port-to-category overrides.
func New(custom map[int]Category) *Grouper {
	if custom == nil {
		custom = make(map[int]Category)
	}
	return &Grouper{custom: custom}
}

// Categorize returns the Category for the given port number.
func (g *Grouper) Categorize(port int) Category {
	if c, ok := g.custom[port]; ok {
		return c
	}
	if c, ok := defaultCategories[port]; ok {
		return c
	}
	return CategoryOther
}

// GroupPorts partitions a slice of port numbers into a map keyed by Category.
func (g *Grouper) GroupPorts(ports []int) map[Category][]int {
	result := make(map[Category][]int)
	for _, p := range ports {
		c := g.Categorize(p)
		result[c] = append(result[c], p)
	}
	return result
}

// Label returns a human-readable string for the port and its category.
func (g *Grouper) Label(port int) string {
	return fmt.Sprintf("%d (%s)", port, g.Categorize(port))
}
