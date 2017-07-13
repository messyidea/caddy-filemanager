package handlers

import (
	"net/http"
	"strings"

	"github.com/messyidea/caddy-filemanager/config"
	"github.com/messyidea/caddy-filemanager/file"
	"github.com/messyidea/caddy-filemanager/page"
	"github.com/messyidea/caddy-filemanager/utils/errors"
)

// ServeSingle serves a single file in an editor (if it is editable), shows the
// plain file, or downloads it if it can't be shown.
func ServeSingle(w http.ResponseWriter, r *http.Request, c *config.Config, u *config.User, i *file.Info) (int, error) {
	var err error

	if err = i.RetrieveFileType(); err != nil {
		return errors.ErrorToHTTPCode(err, true), err
	}

	p := &page.Page{
		Info: &page.Info{
			Name:   i.Name,
			Path:   i.VirtualPath,
			IsDir:  false,
			Data:   i,
			User:   u,
			Config: c,
		},
	}

	// If the request accepts JSON, we send the file information.
	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		return p.PrintAsJSON(w)
	}

	if i.Type == "text" {
		if err = i.Read(); err != nil {
			return errors.ErrorToHTTPCode(err, true), err
		}
	}

	if i.CanBeEdited() && u.AllowEdit {
		p.Data, err = GetEditor(r, i)
		p.Editor = true
		if err != nil {
			return http.StatusInternalServerError, err
		}

		return p.PrintAsHTML(w, "frontmatter", "editor")
	}

	return p.PrintAsHTML(w, "single")
}
