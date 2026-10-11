export const templateRegistry = {
  "gas-grill": {
    id: "gas-grill",
    name: "Gas Grill",
    category: "outdoor-cooking",
    status: "active",
  },
} as const;

export type TemplateId = keyof typeof templateRegistry;
export const DEFAULT_TEMPLATE: TemplateId = "gas-grill";
