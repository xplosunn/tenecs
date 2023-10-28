package testcode

const Comment TestCodeCategory = "comment"

var CommentEverywhere = Create(Comment, "CommentEverywhere", `
// pre-package comment 1
/* pre-package comment 
2 */
package /* post-package */ main

`)
