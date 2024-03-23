package middleware

import "github.com/rulanugrh/tokoku/product/internal/entity/web"

func ValidateRole(tkn string) error {
	claim, err := CheckToken(tkn)
	if err != nil {
		return err
	}
	
	if claim.RoleID != 1 {
		return web.Forbidden("sorry you not admin")
	} else if claim.RoleID != 2 {
		return web.Forbidden("sorry you not owner")
	} else {
		return nil
	}
}
