package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/cheggaaa/pb"
	"github.com/fatih/color"
	"github.com/howeyc/gopass"
	"github.com/joho/godotenv"
	"github.com/ne0z/go-epub"
)

var usr, _ = user.Current()
var tmpfiles []string

// Credential used to store user and password
type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// TokenData used to store JWT token and refresh token
type TokenData struct {
	Refresh string `json:"refresh"`
	Access  string `json:"access"`
}

// Token used to collecting token data
type Token struct {
	Data TokenData `json:"data"`
}

// SearchResult is a struct that used to store received search result
type SearchResult struct {
	Results []struct {
		Hits []struct {
			PrintIsbn13                                string   `json:"print_isbn13"`
			Title                                      string   `json:"title"`
			ProductType                                string   `json:"product_type"`
			PublishedOn                                string   `json:"published_on"`
			EarlyAccess                                int      `json:"early_access"`
			Imprint                                    string   `json:"imprint"`
			AmazonKeywords                             string   `json:"amazon_keywords"`
			Description                                string   `json:"description"`
			Category                                   string   `json:"category"`
			Concept                                    string   `json:"concept"`
			PrimaryLanguage                            string   `json:"primary_language"`
			PrimaryLanguageVersion                     string   `json:"primary_language_version"`
			PrimaryTool                                string   `json:"primary_tool"`
			PrimaryToolVersion                         string   `json:"primary_tool_version"`
			PrimaryToolType                            string   `json:"primary_tool_type"`
			PrimaryToolLanguage                        string   `json:"primary_tool_language"`
			SecondaryTool                              string   `json:"secondary_tool"`
			SecondaryToolVersion                       string   `json:"secondary_tool_version"`
			SecondaryToolType                          string   `json:"secondary_tool_type"`
			SecondaryToolLanguage                      string   `json:"secondary_tool_language"`
			LanguageAudienceKnowledgePrerequisite      string   `json:"language_audience_knowledge_prerequisite"`
			PrimaryToolAudienceKnowledgePrerequisite   string   `json:"primary_tool_audience_knowledge_prerequisite"`
			SecondaryToolAudienceKnowledgePrerequisite string   `json:"secondary_tool_audience_knowledge_prerequisite"`
			Vendor                                     string   `json:"vendor"`
			PublishedYear                              string   `json:"published_year"`
			LatestPacktRank                            int      `json:"latest_packt_rank"`
			Author                                     []string `json:"author"`
			SeoURL                                     string   `json:"seoUrl"`
			MaptSeoURL                                 string   `json:"maptSeoUrl"`
			ImageURL                                   string   `json:"imageUrl"`
			Released                                   string   `json:"released"`
			Length                                     string   `json:"length"`
			ObjectID                                   string   `json:"objectID"`
			HighlightResult                            struct {
				PrintIsbn13 struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"print_isbn13"`
				Title struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"title"`
				ProductType struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"product_type"`
				PublishedOn struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"published_on"`
				EarlyAccess struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"early_access"`
				Imprint struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"imprint"`
				AmazonKeywords struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"amazon_keywords"`
				Description struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"description"`
				Category struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"category"`
				Concept struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"concept"`
				PrimaryLanguage struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"primary_language"`
				PrimaryLanguageVersion struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"primary_language_version"`
				PrimaryTool struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"primary_tool"`
				PrimaryToolVersion struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"primary_tool_version"`
				PrimaryToolType struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"primary_tool_type"`
				PrimaryToolLanguage struct {
					Value            string   `json:"value"`
					MatchLevel       string   `json:"matchLevel"`
					FullyHighlighted bool     `json:"fullyHighlighted"`
					MatchedWords     []string `json:"matchedWords"`
				} `json:"primary_tool_language"`
				SecondaryTool struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"secondary_tool"`
				SecondaryToolVersion struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"secondary_tool_version"`
				SecondaryToolType struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"secondary_tool_type"`
				SecondaryToolLanguage struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"secondary_tool_language"`
				LanguageAudienceKnowledgePrerequisite struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"language_audience_knowledge_prerequisite"`
				PrimaryToolAudienceKnowledgePrerequisite struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"primary_tool_audience_knowledge_prerequisite"`
				SecondaryToolAudienceKnowledgePrerequisite struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"secondary_tool_audience_knowledge_prerequisite"`
				Vendor struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"vendor"`
				PublishedYear struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"published_year"`
				LatestPacktRank struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"latest_packt_rank"`
				Author []struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"author"`
				SeoURL struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"seoUrl"`
				MaptSeoURL struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"maptSeoUrl"`
				ImageURL struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"imageUrl"`
				Released struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"released"`
				Length struct {
					Value        string        `json:"value"`
					MatchLevel   string        `json:"matchLevel"`
					MatchedWords []interface{} `json:"matchedWords"`
				} `json:"length"`
			} `json:"_highlightResult"`
		} `json:"hits"`
		NbHits      int `json:"nbHits"`
		Page        int `json:"page"`
		NbPages     int `json:"nbPages"`
		HitsPerPage int `json:"hitsPerPage"`
		Facets      struct {
			Concept struct {
				Containerization            int `json:"Containerization"`
				DevOps                      int `json:"DevOps"`
				ProgrammingLanguage         int `json:"Programming Language"`
				SystemAdministration        int `json:"System Administration"`
				CloudNative                 int `json:"Cloud Native"`
				Microservices               int `json:"Microservices"`
				CloudComputing              int `json:"Cloud Computing"`
				ApplicationDevelopment      int `json:"Application Development"`
				Serverless                  int `json:"Serverless"`
				WebProgramming              int `json:"Web Programming"`
				ConfigurationManagement     int `json:"Configuration Management"`
				FullStackWebDevelopment     int `json:"Full Stack Web Development"`
				MachineLearning             int `json:"Machine Learning"`
				Blockchain                  int `json:"Blockchain"`
				Concurrency                 int `json:"Concurrency"`
				ContinuousIntegration       int `json:"Continuous Integration"`
				DesignPatterns              int `json:"Design Patterns"`
				FrontEndWebDevelopment      int `json:"Front End Web Development"`
				GUIApplicationDevelopment   int `json:"GUI Application Development"`
				HighPerformanceProgramming  int `json:"High Performance Programming"`
				RESTAPI                     int `json:"REST API"`
				SystemProgramming           int `json:"System Programming"`
				WebServices                 int `json:"Web Services"`
				ApplicationMonitoring       int `json:"Application Monitoring"`
				AugmentedReality            int `json:"Augmented Reality"`
				BuildAutomation             int `json:"Build Automation"`
				CloudSecurity               int `json:"Cloud Security"`
				ContinuousDelivery          int `json:"Continuous Delivery"`
				Cryptocurrency              int `json:"Cryptocurrency"`
				DataAnalysis                int `json:"Data Analysis"`
				DataStructuresAndAlgorithms int `json:"Data Structures and Algorithms"`
				DeepLearning                int `json:"Deep Learning"`
				DistributedComputing        int `json:"Distributed Computing"`
				FunctionalProgramming       int `json:"Functional Programming"`
				GeospatialAnalysis          int `json:"Geospatial Analysis"`
				HomeAutomation              int `json:"Home Automation"`
				ITCertification             int `json:"IT Certification"`
				InformationManagement       int `json:"Information Management"`
				InfrastructureManagement    int `json:"Infrastructure Management"`
				Middleware                  int `json:"Middleware"`
				NetworkSecurity             int `json:"Network Security"`
				Networking                  int `json:"Networking"`
				ProcessManagement           int `json:"Process Management"`
				ProtocolBuffer              int `json:"Protocol Buffer"`
				Servers                     int `json:"Servers"`
				ShellScripting              int `json:"Shell Scripting"`
				SoftwareArchitecture        int `json:"Software Architecture"`
				Virtualization              int `json:"Virtualization"`
			} `json:"concept"`
			Category struct {
				CloudNetworking int `json:"Cloud & Networking"`
				Programming     int `json:"Programming"`
				WebDevelopment  int `json:"Web Development"`
				Data            int `json:"Data"`
				BusinessOther   int `json:"Business & Other"`
				IoTHardware     int `json:"IoT & Hardware"`
				Mobile          int `json:"Mobile"`
				Security        int `json:"Security"`
			} `json:"category"`
			PrimaryTool struct {
				Docker              int `json:"Docker"`
				Kubernetes          int `json:"Kubernetes"`
				Openshift           int `json:"Openshift"`
				Ansible             int `json:"Ansible"`
				CoreOS              int `json:"CoreOS"`
				Jenkins             int `json:"Jenkins"`
				Linux               int `json:"Linux"`
				Terraform           int `json:"Terraform"`
				AWSLambda           int `json:"AWS Lambda"`
				IntelliJIdea        int `json:"IntelliJ Idea"`
				JavaEE              int `json:"Java EE"`
				OpenStack           int `json:"OpenStack"`
				Vagrant             int `json:"Vagrant"`
				GRPC                int `json:"gRPC"`
				AWS                 int `json:"AWS"`
				Angular             int `json:"Angular"`
				ApacheCamel         int `json:"Apache Camel"`
				ArcGIS              int `json:"ArcGIS"`
				Blockchain          int `json:"Blockchain"`
				Buffalo             int `json:"Buffalo"`
				Django              int `json:"Django"`
				Echo                int `json:"Echo"`
				Ethereum            int `json:"Ethereum"`
				Excel               int `json:"Excel"`
				GoogleClassroom     int `json:"Google Classroom"`
				GoogleCloudPlatform int `json:"Google Cloud Platform"`
				GopherJS            int `json:"GopherJS"`
				Hyperledger         int `json:"Hyperledger"`
				MATLAB              int `json:"MATLAB"`
				Puppet              int `json:"Puppet"`
				Quarkus             int `json:"Quarkus"`
				React               int `json:"React"`
				SpringBoot          int `json:"Spring Boot"`
				Unity               int `json:"Unity"`
				WildFly             int `json:"WildFly"`
				Golearn             int `json:"golearn"`
			} `json:"primary_tool"`
			ProductType struct {
				Book         int `json:"Book"`
				Video        int `json:"Video"`
				LearningPath int `json:"Learning Path"`
			} `json:"product_type"`
			PublishedYear struct {
				Num2013 int `json:"2013"`
				Num2014 int `json:"2014"`
				Num2015 int `json:"2015"`
				Num2016 int `json:"2016"`
				Num2017 int `json:"2017"`
				Num2018 int `json:"2018"`
				Num2019 int `json:"2019"`
				Num2020 int `json:"2020"`
				Num2021 int `json:"2021"`
			} `json:"published_year"`
		} `json:"facets"`
		FacetsStats struct {
			PublishedYear struct {
				Min int `json:"min"`
				Max int `json:"max"`
				Avg int `json:"avg"`
				Sum int `json:"sum"`
			} `json:"published_year"`
		} `json:"facets_stats"`
		ExhaustiveFacetsCount bool   `json:"exhaustiveFacetsCount"`
		ExhaustiveNbHits      bool   `json:"exhaustiveNbHits"`
		Query                 string `json:"query"`
		Params                string `json:"params"`
		Index                 string `json:"index"`
		ProcessingTimeMS      int    `json:"processingTimeMS"`
	} `json:"results"`
}

type ArrSearchRequest struct {
	IndexName string `json:"indexName"`
	Params    string `json:"params"`
}

type SearchRequest struct {
	Requests []ArrSearchRequest `json:"requests"`
}

func Login(username string, password string) Token {
	cred := &Credential{
		Username: username,
		Password: password,
	}
	b, _ := json.Marshal(cred)
	req, err := http.NewRequest(
		"POST",
		"https://services.packtpub.com/auth-v1/users/tokens",
		bytes.NewBuffer(b))

	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"88\", \"Google Chrome\";v=\"88\", \";Not A Brand\";v=\"99\"")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://subscription.packtpub.com")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://subscription.packtpub.com/")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var token Token
	json.Unmarshal(body, &token)
	return token
}

func Search(keyword string) SearchResult {
	var result SearchResult

	request := &SearchRequest{
		Requests: []ArrSearchRequest{
			{
				IndexName: "subs_search_prod",
				Params:    "query=" + url.QueryEscape(keyword) + "&hitsPerPage=100&maxValuesPerFacet=100&page=0&attributesToHighlight=%5B%22*%22%5D&filters=product_type%3A%22Video%22%20OR%20product_type%3A%22Book%22%20OR%20product_type%3A%22Project%22%20OR%20product_type%3A%22Learning%20Path%22&ruleContexts=%5B%22subscriptions%22%5D&facets=%5B%22product_type%22%2C%22category%22%2C%22imprint%22%2C%22published_year%22%2C%22concept%22%2C%22primary_tool%22%2C%22primary_language%22%5D&tagFilters=",
			},
		},
	}

	b, _ := json.Marshal(request)
	req, err := http.NewRequest(
		"POST",
		"https://vivzzxfqg1-dsn.algolia.net/1/indexes/*/queries?x-algolia-application-id=VIVZZXFQG1&x-algolia-api-key=945b46c99f0be80981c40d1fb3c7db74",
		bytes.NewBuffer(b))

	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"88\", \"Google Chrome\";v=\"88\", \";Not A Brand\";v=\"99\"")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://subscription.packtpub.com")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://subscription.packtpub.com/")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	return result
}

type Summary struct {
	Title     string `json:"title"`
	ProductID string `json:"productId"`
	Isbn10    string `json:"isbn10"`
	Isbn13    string `json:"isbn13"`
	Isbns     struct {
		Print string `json:"print"`
		Ebook string `json:"ebook"`
	} `json:"isbns"`
	Length          string    `json:"length"`
	PublicationDate time.Time `json:"publicationDate"`
	Authors         []string  `json:"authors"`
	Type            string    `json:"type"`
	OneLiner        string    `json:"oneLiner"`
	About           string    `json:"about"`
	Learn           string    `json:"learn"`
	Features        string    `json:"features"`
	Category        string    `json:"category"`
	Available       bool      `json:"available"`
	Releasing       bool      `json:"releasing"`
	EarlyAccess     bool      `json:"earlyAccess"`
	Pages           int       `json:"pages"`
	ShopURL         string    `json:"shopUrl"`
	ReadURL         string    `json:"readUrl"`
	Meta            struct {
		Category struct {
			CategoryName string `json:"category_name"`
		} `json:"category"`
		Concepts struct {
			ConceptName string `json:"concept_name"`
		} `json:"concepts"`
		Language struct {
			LanguageName string `json:"language_name"`
		} `json:"language"`
		LanguageVersion struct {
			LanguageVersionName string `json:"language_version_name"`
		} `json:"languageVersion"`
		Tool struct {
			ToolName string `json:"tool_name"`
		} `json:"tool"`
		Vendor struct {
			Vendor string `json:"vendor"`
		} `json:"vendor"`
	} `json:"meta"`
	InStore                  bool          `json:"inStore"`
	CoverImage               string        `json:"coverImage"`
	InSubs                   bool          `json:"inSubs"`
	Licensed                 bool          `json:"licensed"`
	Publisher                string        `json:"publisher"`
	DistributionRestrictions []interface{} `json:"distributionRestrictions"`
}

func GetSummary(isbn string) Summary {
	var summary Summary

	req, err := http.NewRequest(
		"GET",
		"https://static.packt-cdn.com/products/"+isbn+"/summary",
		nil,
	)

	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"88\", \"Google Chrome\";v=\"88\", \";Not A Brand\";v=\"99\"")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://subscription.packtpub.com")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://subscription.packtpub.com/")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &summary)
	return summary
}

type TOC struct {
	ProductID string `json:"productId"`
	Prefaces  []struct {
		ID       string `json:"id"`
		Title    string `json:"title"`
		Sections []struct {
			ID          string `json:"id"`
			Title       string `json:"title"`
			ContentType string `json:"contentType"`
		} `json:"sections"`
	} `json:"prefaces"`
	Appendices []struct {
		ID       string `json:"id"`
		Title    string `json:"title"`
		Sections []struct {
			ID          string `json:"id"`
			Title       string `json:"title"`
			ContentType string `json:"contentType"`
		} `json:"sections"`
	} `json:"appendices"`
	Chapters []struct {
		ID       string `json:"id"`
		Title    string `json:"title"`
		Sections []struct {
			ID          string `json:"id"`
			Title       string `json:"title"`
			ContentType string `json:"contentType"`
		} `json:"sections"`
	} `json:"chapters"`
}

func GetToc(isbn string) TOC {
	var toc TOC

	req, err := http.NewRequest(
		"GET",
		"https://static.packt-cdn.com/products/"+isbn+"/toc",
		nil,
	)

	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"88\", \"Google Chrome\";v=\"88\", \";Not A Brand\";v=\"99\"")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://subscription.packtpub.com")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://subscription.packtpub.com/")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &toc)
	return toc
}

type Author struct {
	ID           string   `json:"id"`
	EpicAuthorID string   `json:"epicAuthorId"`
	Author       string   `json:"author"`
	Description  string   `json:"description"`
	Products     []string `json:"products"`
	URLKey       string   `json:"urlKey"`
}

func GetAuthor(authorId string) Author {
	var author Author

	req, err := http.NewRequest(
		"GET",
		"https://static.packt-cdn.com/authors/"+authorId,
		nil,
	)

	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"88\", \"Google Chrome\";v=\"88\", \";Not A Brand\";v=\"99\"")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://subscription.packtpub.com")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://subscription.packtpub.com/")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &author)
	return author
}

type Page struct {
	Data string `json:"data"`
}

func GetPage(token string, isbn string, chapterId string, pageId string) Page {
	var page Page

	req, err := http.NewRequest(
		"GET",
		"https://services.packtpub.com/products-v1/products/"+isbn+"/"+chapterId+"/"+pageId,
		nil,
	)

	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"88\", \"Google Chrome\";v=\"88\", \";Not A Brand\";v=\"99\"")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://subscription.packtpub.com")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://subscription.packtpub.com/")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &page)
	return page
}

func DownloadPage(url string) (string, error) {
	req, err := http.NewRequest(
		"GET",
		url,
		nil,
	)

	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"88\", \"Google Chrome\";v=\"88\", \";Not A Brand\";v=\"99\"")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Origin", "https://subscription.packtpub.com")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://subscription.packtpub.com/")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func RemoveHtmlTag(in string) string {
	// regex to match html tag
	const pattern = `(<h1).*(</h1>)`
	r := regexp.MustCompile(pattern)
	groups := r.FindAllString(in, -1)
	// should replace long string first
	sort.Slice(groups, func(i, j int) bool {
		return len(groups[i]) > len(groups[j])
	})
	for _, group := range groups {
		if strings.TrimSpace(group) != "" {
			in = strings.ReplaceAll(in, group, "")
		}
	}
	return in
}

func fixImgSrc(in string, isbn string) string {
	var re = regexp.MustCompile(`(<img src=").*(/graphics)`)
	s := re.ReplaceAllString(in, `${1}https://static.packt-cdn.com/products/`+isbn+`${2}`)
	return s
}

func DownloadAsMobi(token string, isbn string) string {
	fileName := DownloadAsEpub(token, isbn)

	_, err := os.Stat("/Applications/calibre.app/Contents/MacOS//ebook-convert")
	app := ""
	if os.IsNotExist(err) {
		app = "/usr/bin/ebook-convert"
	} else {
		app = "/Applications/calibre.app/Contents/MacOS//ebook-convert"
	}

	cmd := &exec.Cmd{
		Path: app,
		Args: []string{app, fileName + ".epub", fileName + ".mobi"},
	}
	cmd.Start()
	cmd.Wait()
	os.Remove(fileName + ".epub")
	return fileName
}

func fixImgInternalSrc(in string, isbn string) (string, [][]string) {
	re := regexp.MustCompile("<img src=\"(.*?)\"")
	matches := re.FindAllStringSubmatch(in, -1)
	if len(matches) == 0 {
		return in, nil
	}

	re = regexp.MustCompile(`/graphics/(.*?)/graphics/(.*?)/`)
	s := re.ReplaceAllString(in, `../images/`)

	re = regexp.MustCompile(`(<img)(.*?)(>)`)
	s = re.ReplaceAllString(s, `${1}${2}/${3}`)

	re = regexp.MustCompile(`(<br)(.*?)(>)`)
	s = re.ReplaceAllString(s, `${1}${2}/${3}`)
	return s, matches
}

func FixImgUrlError(url string) string {
	re := regexp.MustCompile(`/graphics/(.*?)/`)
	s := re.ReplaceAllString(url, ``)
	return s
}

func EmbedEPubImage(e *epub.Epub, isbn string, urls [][]string) {
	for _, imgURL := range urls {
		imagePath := usr.HomeDir + "/.packt_tmp/" + filepath.Base(imgURL[1])
		_, err := grab.Get(imagePath,
			"https://static.packt-cdn.com/products/"+isbn+"/"+FixImgUrlError(imgURL[1]))
		if err != nil {
			color.Red("Error occured!\r\n" + err.Error())
			color.Red("https://static.packt-cdn.com/products/" + isbn + "/" + FixImgUrlError(imgURL[1]))
			color.Red(imgURL[1])
			continue
		}
		e.AddImage(imagePath, "../images/"+filepath.Base(imgURL[1]))
		tmpfiles = append(tmpfiles, imagePath)
	}
}

func cleanUpFiles() {
	dirname := usr.HomeDir + "/.packt_tmp/"

	d, err := os.Open(dirname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == ".png" || filepath.Ext(file.Name()) == ".jpg" || filepath.Ext(file.Name()) == "jpeg" {
				os.Remove(file.Name())
			}
		}
	}
}

func removeTmpFiles() {
	for _, file := range tmpfiles {
		os.Remove(file)
	}
}

func DownloadAsEpub(token string, isbn string) string {
	summary := GetSummary(isbn)
	author := GetAuthor(summary.Authors[0])
	toc := GetToc(isbn)
	e := epub.NewEpub(summary.Title)
	e.SetAuthor(author.Author)
	e.SetDescription(summary.About)
	e.SetIdentifier(isbn)

	coverPath := usr.HomeDir + "/.packt_tmp/" + filepath.Base(summary.CoverImage)

	_, err := grab.Get(coverPath,
		summary.CoverImage)
	if err != nil {
		log.Fatal(err)
	}
	tmpfiles = append(tmpfiles, coverPath)
	coverImagePath, _ := e.AddImage(coverPath, filepath.Base(summary.CoverImage))
	e.SetCover(coverImagePath, "")

	bar := pb.StartNew(len(toc.Chapters))
	curr := 0
	for _, chapters := range toc.Chapters {
		bar.Increment()
		pageURL := GetPage(token, isbn, chapters.ID, chapters.Sections[0].ID)
		downloadedPage, err := DownloadPage(pageURL.Data)
		if err != nil {
			log.Fatal(pageURL)
		}
		pageData, imgURL := fixImgInternalSrc(downloadedPage, isbn)
		if len(imgURL) > 0 {
			EmbedEPubImage(e, isbn, imgURL)
		}
		e.AddSection(pageData, chapters.Title, "", "")
		for _, subchapter := range chapters.Sections {
			if curr > 0 {
				pageURL := GetPage(token, isbn, chapters.ID, subchapter.ID)
				downloadedPage, err := DownloadPage(pageURL.Data)
				if err != nil {
					log.Fatal(pageURL)
				}
				pageData, imgURL := fixImgInternalSrc(downloadedPage, isbn)
				if len(imgURL) > 0 {
					EmbedEPubImage(e, isbn, imgURL)
				}
				e.AddSection(pageData, subchapter.Title, "", "")
			}
			curr++
		}
		curr = 0
	}
	bar.Finish()

	err = e.Write(summary.Title + ".epub")
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}
	return summary.Title
}

type RefreshCredential struct {
	Refresh string
}

func RefreshToken() Token {
	cred := &RefreshCredential{
		Refresh: os.Getenv("TOKEN"),
	}
	b, _ := json.Marshal(cred)
	req, err := http.NewRequest(
		"POST",
		"https://services.packtpub.com/auth-v1/users/tokens",
		bytes.NewBuffer(b))

	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"88\", \"Google Chrome\";v=\"88\", \";Not A Brand\";v=\"99\"")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://subscription.packtpub.com")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://subscription.packtpub.com/")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var token Token
	json.Unmarshal(body, &token)
	return token
}

func main() {
	cleanUpFiles()
	if len(os.Args) == 2 && os.Args[1] == "login" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Username : ")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)
		fmt.Print("Password : ")
		password, _ := gopass.GetPasswdMasked()
		resp := Login(username, string(password))

		err := ioutil.WriteFile(
			usr.HomeDir+"/.packt_config",
			[]byte("TOKEN="+resp.Data.Access+"\r\nREFRESH="+resp.Data.Refresh),
			0644)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	if len(os.Args) < 3 {
		color.Red("Command error!\n")
		fmt.Println("$ " + os.Args[0] + " <options> <arguments>\n")
		fmt.Println("Available options:\n- login\n- search <keyword>\n- mobi <isbn>\n- epub <isbn>")
		return
	}

	if _, err := os.Stat(usr.HomeDir + "/.packt_tmp"); err == nil {
		os.Mkdir(usr.HomeDir+"/.packt_tmp", 0744)
	}

	if _, err := os.Stat(usr.HomeDir + "/.packt_config"); err == nil {
		godotenv.Load(usr.HomeDir + "/.packt_config")
	} else {
		color.Red("Please login first!")
		os.Exit(1)
	}

	oToken := RefreshToken()
	if len(oToken.Data.Access) > 0 {
		ioutil.WriteFile(
			usr.HomeDir+"/.packt_config",
			[]byte("TOKEN="+oToken.Data.Access+"\r\nREFRESH="+oToken.Data.Refresh),
			0644)
	}

	switch opt := os.Args[1]; opt {

	case "search":
		searchResult := Search(os.Args[2])
		color.Blue("Results:")
		for _, item := range searchResult.Results[0].Hits {
			fmt.Println(item.PrintIsbn13 + " - " + item.PublishedYear + " - " + item.Title)
		}
	case "mobi":
		fileName := DownloadAsMobi(oToken.Data.Access, os.Args[2])
		fmt.Println("Output : " + fileName + ".mobi")
	case "epub":
		fileName := DownloadAsEpub(oToken.Data.Access, os.Args[2])
		fmt.Println("Output : " + fileName + ".epub")
	}
	removeTmpFiles()
}
