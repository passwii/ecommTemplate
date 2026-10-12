export type CompanyStory = {
  heading: string;
  body: string;
  imageUrl: string;
  imageAlt: string;
};

export type CompanyCredential = {
  name: string;
  issuer: string;
  details: string;
};

export type CompanyCaseStudy = {
  customer: string;
  title: string;
  details: string;
};

export type CompanyFaq = {
  question: string;
  answer: string;
};

export type CompanyContent = {
  about: CompanyStory;
  manufacturing: {
    heading: string;
    overview: string;
    qualityHeading: string;
    qualityDetails: string;
  };
  credentials: CompanyCredential[];
  caseStudies: CompanyCaseStudy[];
  faqs: CompanyFaq[];
};

const STORAGE_KEY = "ember-field.storefront.company-content.v1";

export function createEmptyCompanyContent(): CompanyContent {
  return {
    about: { heading: "", body: "", imageUrl: "", imageAlt: "" },
    manufacturing: { heading: "", overview: "", qualityHeading: "", qualityDetails: "" },
    credentials: [],
    caseStudies: [],
    faqs: [],
  };
}

export function getStoredCompanyContent(): CompanyContent {
  if (typeof window === "undefined") return createEmptyCompanyContent();

  try {
    const stored = window.localStorage.getItem(STORAGE_KEY);
    if (!stored) return createEmptyCompanyContent();

    const parsed = JSON.parse(stored) as Partial<CompanyContent>;
    const empty = createEmptyCompanyContent();
    return {
      about: { ...empty.about, ...parsed.about },
      manufacturing: { ...empty.manufacturing, ...parsed.manufacturing },
      credentials: Array.isArray(parsed.credentials) ? parsed.credentials : [],
      caseStudies: Array.isArray(parsed.caseStudies) ? parsed.caseStudies : [],
      faqs: Array.isArray(parsed.faqs) ? parsed.faqs : [],
    };
  } catch {
    return createEmptyCompanyContent();
  }
}

export function saveStoredCompanyContent(content: CompanyContent): void {
  window.localStorage.setItem(STORAGE_KEY, JSON.stringify(content));
}
