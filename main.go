package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type Hotel struct {
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	Location          string  `json:"location"`
	Rating            float64 `json:"rating"`
	Contact           string  `json:"contact"`
	Reviews           string  `json:"reviews"`
	PropertyAmenities string  `json:"propertyAmenities"`
	RoomFeatures      string  `json:"roomFeatures"`
}

const (
	host   = "www.tripadvisor.com"
	scheme = "https"
)

var (
	baseUrl = &url.URL{
		Host:   host,
		Scheme: scheme,
	}
)

func main() {
	var propertyAmenities string
	var roomFeatures string

	browser := rod.New().MustConnect()
	page := browser.
		Timeout(time.Second * 60).
		MustPage(baseUrl.String()).MustWindowFullscreen()

	// Taking screenshot of website
	page.MustWaitLoad().MustScreenshot("tripadvisor.png")

	// Finding search input
	searchEls, err := page.
		Timeout(time.Second * 10).
		Elements(`input[type="search"]`)
	handleError(err)

	if len(searchEls) > 1 {
		// Taking second index of search input to search hotel
		searchEls[1].Input("The Hay-Adams")
	} else {
		// Taking first index of search input to search hotel
		searchEls[0].Input("The Hay-Adams")
	}

	// Opening hotel page after selecting in dropdown list
	hotelNameButtonEl, err := page.
		Timeout(time.Second * 10).
		Element(`div[class="EtzER"]`)
	handleError(err)
	hotelNameButtonEl.Click(proto.InputMouseButtonLeft)

	// Getting hotel name element
	hotelNameEl, err := page.
		Timeout(time.Second * 10).
		Element(`h1[class="QdLfr b d Pn"]`)
	handleError(err)

	// Getting hotel name
	hotelName, err := hotelNameEl.
		Timeout(time.Second * 10).
		Text()
	handleError(err)

	// Getting hotel description element
	descriptionEl, err := page.
		Timeout(time.Second * 10).
		Element(`div[class="fIrGe _T"]`)
	handleError(err)

	// Getting hotel description
	hotelDescription, err := descriptionEl.
		Timeout(time.Second * 10).
		Text()
	handleError(err)

	// Getting hotel reviews element
	hotelReviewsEl, err := page.
		Timeout(time.Second * 10).
		Element(`span[class="qqniT"]`)
	handleError(err)

	// Getting hotel reviews
	hotelReviews, err := hotelReviewsEl.
		Timeout(time.Second * 10).
		Text()
	handleError(err)

	// Getting hotel contact element
	hotelContactEl, err := page.
		Timeout(time.Second * 10).
		Element(`span[class="zNXea NXOxh NjUDn"]`)
	handleError(err)

	// Getting hotel contact
	hotelContact, err := hotelContactEl.
		Timeout(time.Second * 10).
		Text()
	handleError(err)

	// Getting hotel location element
	hotelLocationEl, err := page.
		Timeout(time.Second * 10).
		Element(`span[class="fHvkI PTrfg"]`)
	handleError(err)

	// Getting hotel location
	hotelLocation, err := hotelLocationEl.
		Timeout(time.Second * 10).
		Text()
	handleError(err)

	// Getting hotel rating element
	hotelRatingEl, err := page.
		Timeout(time.Second * 10).
		Element(`span[class="IHSLZ P"]`)
	handleError(err)

	// Getting hotel rating
	rating, err := hotelRatingEl.
		Timeout(time.Second * 10).
		Text()
	handleError(err)

	// Converting rating to float from string data type
	hotelRating, err := strconv.ParseFloat(rating, 64)
	if err != nil {
		fmt.Println(err)
	}

	// Getting property amenities and room features elements
	PropertyAmenitiesAndRoomFeaturesEls, err := page.
		Timeout(time.Second * 10).
		Elements(`div[class="OsCbb K"]`)
	handleError(err)

	if len(PropertyAmenitiesAndRoomFeaturesEls) > 1 {
		// Taking first index of PropertyAmenitiesAndRoomFeaturesEls to get property amenities
		propertyAmenities, err = PropertyAmenitiesAndRoomFeaturesEls[0].
			Timeout(time.Second * 10).
			Text()
		handleError(err)

		// Taking second index of PropertyAmenitiesAndRoomFeaturesEls to get room features
		roomFeatures, err = PropertyAmenitiesAndRoomFeaturesEls[1].
			Timeout(time.Second * 10).
			Text()
		handleError(err)
	}

	// Using marshal function to convert object to JSON
	hotel, err := json.Marshal(
		Hotel{
			Name:              hotelName,
			Description:       hotelDescription,
			Location:          hotelLocation,
			Rating:            hotelRating,
			Contact:           hotelContact,
			Reviews:           hotelReviews,
			PropertyAmenities: propertyAmenities,
			RoomFeatures:      roomFeatures,
		})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print("hotel:", string(hotel))
	time.Sleep(5 * time.Minute)
}

func handleError(err error) {
	var evalErr *rod.ErrEval
	if errors.Is(err, context.DeadlineExceeded) {
		fmt.Println("timeout err")
	} else if errors.As(err, &evalErr) {
		fmt.Println(evalErr.LineNumber)
	} else if err != nil {
		fmt.Println("can't handle", err)
	}
}
