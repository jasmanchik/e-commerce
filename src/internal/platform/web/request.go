package web

import (
	"encoding/json"
	"github.com/go-playground/locales/en"
	"github.com/pkg/errors"
	"net/http"
	"reflect"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_tranlations "gopkg.in/go-playground/validator.v9/translations/en"
)

var (
	validate   = validator.New()
	translator *ut.UniversalTranslator
)

func init() {

	enLocate := en.New()

	translator = ut.New(enLocate, enLocate)

	lang, _ := translator.GetTranslator("en")

	err := en_tranlations.RegisterDefaultTranslations(validate, lang)
	if err != nil {
		errors.Wrap(err, "translation")
	}

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func Decode(r *http.Request, val interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(val); err != nil {
		return NewRequestError(err, http.StatusBadRequest)
	}

	if err := validate.Struct(val); err != nil {
		vErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		lang, _ := translator.GetTranslator("en")
		var fields []FieldError
		for _, verror := range vErrors {
			field := FieldError{
				Field: verror.Field(),
				Error: verror.Translate(lang),
			}
			fields = append(fields, field)
		}

		return &Error{
			Err:    errors.New("field validation error"),
			Status: http.StatusBadRequest,
			Fields: fields,
		}
	}

	return nil
}
