package v1

import (
    "net/http"
)
type MotdAPI struct {}

func NewMotdAPI() *MotdAPI {
    return &MotdAPI{}
}


func (api *MotdAPI) Get(w http.ResponseWriter, r *http.Request) {
    message := `
    Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis lorem elit, consectetur ut posuere quis, rhoncus at lacus. Mauris sem ipsum, pulvinar et rhoncus ut, vestibulum a ipsum. Aenean sit amet leo ornare, maximus elit porttitor, tincidunt dolor. Vestibulum eu ipsum eu sem pellentesque tincidunt. Donec eget nisl rutrum, tincidunt nisi et, egestas est. Sed maximus pulvinar viverra. Nullam tincidunt, justo quis efficitur gravida, erat magna porttitor turpis, sit amet dictum metus risus sit amet odio. Duis rhoncus rhoncus dolor, euismod pharetra metus fringilla tristique. Aliquam sodales, sapien vel viverra ullamcorper, ex ante feugiat tortor, ac viverra nibh neque sit amet turpis. Etiam vitae ante faucibus, convallis erat id, tristique augue. Suspendisse efficitur odio sed sapien bibendum, et laoreet est luctus. Ut pretium vitae mi a eleifend. Aliquam rutrum, dolor sit amet ultricies pellentesque, nisl sem vehicula magna, at auctor ex magna ut purus. Sed dapibus, augue quis lacinia fringilla, enim libero tincidunt erat, eu aliquet lacus mauris vitae orci.

Pellentesque ut condimentum ex. Suspendisse gravida magna felis, ac efficitur lacus tristique quis. Etiam faucibus porttitor turpis, quis efficitur orci laoreet ac. Cras et rhoncus nisl, quis elementum libero. Morbi aliquet lacus mi, eu ullamcorper libero luctus quis. Nulla interdum pulvinar semper. Nulla ullamcorper nisi ut facilisis eleifend. Mauris in nunc aliquet, cursus enim eu, congue urna. Integer tempus dolor egestas arcu varius faucibus. Maecenas vulputate viverra dolor, eget imperdiet nunc feugiat eget.

Donec lacinia felis et eros blandit ornare. Etiam maximus massa ante, eu mattis magna vehicula lobortis. Donec ullamcorper libero et massa auctor rutrum. Sed vulputate bibendum lacus vitae venenatis. Donec malesuada ut velit nec faucibus. Nunc posuere egestas justo, non hendrerit lacus pellentesque a. Cras vitae dapibus leo. Integer tellus justo, aliquet placerat nisi gravida, suscipit viverra velit. Nam vel sapien at ante rhoncus lobortis vel sit amet felis.
`

    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(message))
}

