# fillinform

port from [HTML::FillinForm::Lite](http://search.cpan.org/~gfuji/HTML-FillInForm-Lite-1.13/lib/HTML/FillInForm/Lite.pm)

HTML::FillinForm::Lite is licensed to Goro Fuji.

## installation

    go get -u github.com/sheercat/fillinform

## Usage
This product is alpha version yet.

use html/template

    import (
       "github.com/sheercat/fillinform"
       "html/template"
    )
    
    ...
    
    writer := fillinform.FillWriter(w, fdat, nil)
    html.ExecuteTemplate(writer, "layout", map[string]interface{}{"reqParams": reqParams})

use pongo2

    import (
       "github.com/sheercat/fillinform"
       "github.com/flosch/pongo2"
    )
    
    ....
    
    bytes, err := tpl.ExecuteBytes(pongo2.Context{})
    if err != nil {
       http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    
    bytes, err = fillinform.Fill(&bytes, formData.(map[string]interface{}), nil)
    if err != nil {
       http.Error(w, err.Error(), http.StatusInternalServerError)
    }


## License



fillinform licensed under the MIT



