package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	// serializers
	"github.com/PartyLich/hex-microservice/serializer/json"
	"github.com/PartyLich/hex-microservice/serializer/msgpack"
	// domain logic
	"github.com/PartyLich/hex-microservice/shortUrl"
)

type RedirectHandler interface {
	Get(http.ResponseWriter, *http.Request)
	Post(http.ResponseWriter, *http.Request)
}

type handler struct {
	redirectService shortUrl.RedirectService
}

// NewHandler creates an instance of the handler type
func NewHandler(redirectService shortUrl.RedirectService) RedirectHandler {
	return &handler{redirectService: redirectService}
}

// setupResponse creates appropriate http response header
func setupResponse(res http.ResponseWriter, contentType string, body []byte, statusCode int) {
	res.Header().Set("Content-Type", contentType)
	res.WriteHeader(statusCode)
	_, err := res.Write(body)
	if err != nil {
		log.Println(err)
	}
}

// get requested serializer
func (h *handler) serializer(contentType string) shortUrl.RedirectSerializer {
	switch contentType {
	case "application/x-msgpack":
		return &msgpack.Redirect{}
	default:
		return &json.Redirect{}
	}
}

// redirect to stored url
func (h *handler) Get(res http.ResponseWriter, req *http.Request) {
	code := chi.URLParam(req, "code")
	redirect, err := h.redirectService.Find(code)
	if err != nil {
		if errors.Cause(err) == shortUrl.ErrRedirectNotFound {
			http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	http.Redirect(res, req, redirect.URL, http.StatusMovedPermanently)
}

// create code for url redirect
func (h *handler) Post(res http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	redirect, err := h.serializer(contentType).Decode(requestBody)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	err = h.redirectService.Store(redirect)
	if err != nil {
		if errors.Cause(err) == shortUrl.ErrRedirectInvalid {
			http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	responseBody, err := h.serializer(contentType).Encode(redirect)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	setupResponse(res, contentType, responseBody, http.StatusCreated)
}
