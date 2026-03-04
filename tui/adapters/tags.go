package adapters

import "github.com/tonysyu/gqlxp/gql"

// objectTags returns display tags for a GraphQL Object type.
// "Entity" is added when the object has a @key directive (Apollo Federation).
// "Node" is added when the object implements the Node interface.
func objectTags(obj *gql.Object) []string {
	var tags []string
	for _, dir := range obj.Directives() {
		if dir.Name() == "key" {
			tags = append(tags, "Entity")
			break
		}
	}
	for _, iface := range obj.Interfaces() {
		if iface == "Node" {
			tags = append(tags, "Node")
			break
		}
	}
	return tags
}
