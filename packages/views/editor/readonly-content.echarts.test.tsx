import { describe, it, expect, vi } from "vitest";
import { render, cleanup } from "@testing-library/react";
import * as fc from "fast-check";

// ---------------------------------------------------------------------------
// Mocks
// ---------------------------------------------------------------------------

// Mock echarts-for-react to avoid canvas/DOM issues in jsdom
vi.mock("echarts-for-react", () => ({
  default: ({ option, className }: { option: unknown; className?: string }) => (
    <div
      data-testid="echarts-chart"
      data-option={JSON.stringify(option)}
      className={className}
    >
      ECharts Chart
    </div>
  ),
}));

// Mock hooks from @multica/core that ReadonlyContent depends on
vi.mock("@multica/core/paths", () => ({
  useWorkspacePaths: () => ({
    issueDetail: (id: string) => `/issues/${id}`,
  }),
  useWorkspaceSlug: () => "test-workspace",
}));

// Mock navigation hook used by ReadonlyLink
vi.mock("../navigation", () => ({
  useNavigation: () => ({
    push: vi.fn(),
    openInNewTab: vi.fn(),
  }),
}));

// Mock IssueMentionCard to avoid deep dependency chain
vi.mock("../issues/components/issue-mention-card", () => ({
  IssueMentionCard: ({ issueId }: { issueId: string }) => (
    <span data-testid="issue-mention">{issueId}</span>
  ),
}));

// Mock ImageLightbox to avoid complex DOM interactions
vi.mock("./extensions/image-view", () => ({
  ImageLightbox: () => null,
}));

// Mock preprocessMarkdown to pass content through unchanged
vi.mock("./utils/preprocess", () => ({
  preprocessMarkdown: (value: string) => value,
}));

// Mock link-hover-card to avoid floating-ui DOM issues in jsdom
vi.mock("./link-hover-card", () => ({
  useLinkHover: () => ({
    visible: false,
    href: "",
    anchorEl: null,
    cardRef: { current: null },
    onCardEnter: vi.fn(),
    onCardLeave: vi.fn(),
  }),
  LinkHoverCard: () => null,
}));

import { ReadonlyContent } from "./readonly-content";

// ---------------------------------------------------------------------------
// Arbitraries
// ---------------------------------------------------------------------------

/**
 * Generate JSON-serializable objects (not primitives).
 * ECharts options are objects or arrays, so we restrict to those shapes.
 */
const jsonObjectArb = fc.oneof(
  fc.dictionary(fc.string(), fc.jsonValue()),
  fc.array(fc.jsonValue()),
);

/**
 * Generate strings that are NOT valid JSON.
 * Excludes empty strings since an empty echarts block produces no <code> element.
 */
const invalidJsonArb = fc.oneof(
  fc.constant("{broken"),
  fc.constant("undefined"),
  fc.constant("[1, 2,"),
  fc.constant('{"key": }'),
  fc.constant("not json at all"),
  fc.string({ minLength: 1 }).filter((s) => {
    try {
      JSON.parse(s);
      return false; // Valid JSON — exclude it
    } catch {
      return true; // Invalid JSON — include it
    }
  }),
);

// ---------------------------------------------------------------------------
// Property-Based Tests
// ---------------------------------------------------------------------------

describe("ReadonlyContent ECharts rendering (Property-Based Tests)", () => {
  /**
   * Feature: echarts-markdown-renderer, Property 3:
   * Valid JSON echarts blocks render as charts in ReadonlyContent with transparent background
   *
   * For any valid JSON object used as the body of an `echarts` fenced code block
   * rendered by `ReadonlyContent`, the output should contain an `EChartsChart`
   * component with `className="echarts-chart-wrapper"` and the `option` prop
   * should have `backgroundColor` set to `'transparent'` regardless of whether
   * the original option specified a background color.
   *
   * Validates: Requirements 4.1, 4.2, 4.4, 4.5
   */
  it("Property 3: for any fc.jsonValue() object, ReadonlyContent renders EChartsChart with className='echarts-chart-wrapper' and option.backgroundColor === 'transparent'", () => {
    fc.assert(
      fc.property(jsonObjectArb, (jsonObj) => {
        const jsonString = JSON.stringify(jsonObj);
        const markdown = `\`\`\`echarts\n${jsonString}\n\`\`\``;

        const { getByTestId } = render(
          <ReadonlyContent content={markdown} />,
        );

        // Assert EChartsChart is rendered
        const chart = getByTestId("echarts-chart");
        expect(chart).toBeInTheDocument();

        // Assert className is correct (Requirement 4.5)
        expect(chart).toHaveClass("echarts-chart-wrapper");

        // Assert option.backgroundColor is 'transparent' (Requirement 4.4)
        const optionAttr = chart.getAttribute("data-option");
        expect(optionAttr).not.toBeNull();
        const parsedOption = JSON.parse(optionAttr!);
        expect(parsedOption.backgroundColor).toBe("transparent");

        cleanup();
      }),
      { numRuns: 100 },
    );
  });

  /**
   * Feature: echarts-markdown-renderer, Property 4:
   * Invalid JSON echarts blocks fall back to code in ReadonlyContent
   *
   * For any string that is not valid JSON used as the body of an `echarts` fenced
   * code block rendered by `ReadonlyContent`, the output should contain a fallback
   * code block with the original string content unchanged, and no `EChartsChart`
   * component.
   *
   * Validates: Requirements 4.3
   */
  it("Property 4: for any non-JSON string, ReadonlyContent renders a fallback code block containing the original string and no EChartsChart", () => {
    fc.assert(
      fc.property(invalidJsonArb, (invalidJson) => {
        const markdown = `\`\`\`echarts\n${invalidJson}\n\`\`\``;

        const { container, queryByTestId } = render(
          <ReadonlyContent content={markdown} />,
        );

        // Assert NO EChartsChart is rendered (Requirement 4.3)
        const chart = queryByTestId("echarts-chart");
        expect(chart).not.toBeInTheDocument();

        // Assert fallback code block is rendered and contains the original string.
        // react-markdown and lowlight may add/preserve a trailing newline in the
        // rendered code element's textContent, so we compare after stripping it
        // from both sides.
        const codeBlock = container.querySelector("code");
        expect(codeBlock).toBeInTheDocument();
        const renderedText = (codeBlock?.textContent ?? "").replace(/\n$/, "");
        const expectedText = invalidJson.replace(/\n$/, "");
        expect(renderedText).toBe(expectedText);

        cleanup();
      }),
      { numRuns: 100 },
    );
  });
});
