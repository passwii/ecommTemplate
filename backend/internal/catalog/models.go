package catalog

// Collection is a storefront product grouping shown on the home and collection
// pages.
type Collection struct {
	ID          int64  `json:"id"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SortOrder   int    `json:"sortOrder"`
}

// ProductImage is one ordered gallery image for a product.
type ProductImage struct {
	URL       string `json:"url"`
	Alt       string `json:"alt"`
	SortOrder int    `json:"sortOrder"`
}

// ProductSpecification is a name/value specification row.
type ProductSpecification struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ProductResource is a downloadable or linked product resource.
type ProductResource struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	FileType    string `json:"fileType,omitempty"`
	Description string `json:"description,omitempty"`
}

// Product is the public storefront representation of a product. Draft products
// are never returned by the public API.
type Product struct {
	ID                   string                 `json:"id"`
	Slug                 string                 `json:"slug"`
	Name                 string                 `json:"name"`
	Model                string                 `json:"model,omitempty"`
	SKU                  string                 `json:"sku,omitempty"`
	Collection           Collection             `json:"collection"`
	ShortDescription     string                 `json:"shortDescription"`
	Description          string                 `json:"description"`
	Images               []ProductImage         `json:"images"`
	Specifications       []ProductSpecification `json:"specifications"`
	Applications         []string               `json:"applications"`
	Resources            []ProductResource      `json:"resources"`
	PriceAmount          *float64               `json:"priceAmount,omitempty"`
	PriceCurrency        string                 `json:"priceCurrency,omitempty"`
	CompareAtPrice       *float64               `json:"compareAtPrice,omitempty"`
	Badge                string                 `json:"badge,omitempty"`
	Fuel                 string                 `json:"fuel,omitempty"`
	Burners              *int                   `json:"burners,omitempty"`
	CookingArea          string                 `json:"cookingArea,omitempty"`
	Finish               string                 `json:"finish,omitempty"`
	MinimumOrderQuantity string                 `json:"minimumOrderQuantity,omitempty"`
	LeadTime             string                 `json:"leadTime,omitempty"`
	Featured             bool                   `json:"featured"`
	SortOrder            int                    `json:"sortOrder"`
	ContentStatus        string                 `json:"contentStatus"`
}

// CompanyStory is the About section content.
type CompanyStory struct {
	Heading  string `json:"heading"`
	Body     string `json:"body"`
	ImageURL string `json:"imageUrl"`
	ImageAlt string `json:"imageAlt"`
}

// CompanyManufacturing is the manufacturing and quality section content.
type CompanyManufacturing struct {
	Heading        string `json:"heading"`
	Overview       string `json:"overview"`
	QualityHeading string `json:"qualityHeading"`
	QualityDetails string `json:"qualityDetails"`
}

// CompanyCredential is a certification or credential entry. Only entries with
// verifiable evidence may be presented as facts.
type CompanyCredential struct {
	Name    string `json:"name"`
	Issuer  string `json:"issuer"`
	Details string `json:"details"`
}

// CompanyCaseStudy is a customer case study entry.
type CompanyCaseStudy struct {
	Customer string `json:"customer"`
	Title    string `json:"title"`
	Details  string `json:"details"`
}

// CompanyFAQ is a question and answer pair.
type CompanyFAQ struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// CompanyContent is the public company/trust page content.
type CompanyContent struct {
	About         CompanyStory         `json:"about"`
	Manufacturing CompanyManufacturing `json:"manufacturing"`
	Credentials   []CompanyCredential  `json:"credentials"`
	CaseStudies   []CompanyCaseStudy   `json:"caseStudies"`
	FAQs          []CompanyFAQ         `json:"faqs"`
	ContentStatus string               `json:"contentStatus"`
}
