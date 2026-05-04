import { describe, it, expect, vi } from "vitest";
import { render, cleanup } from "@testing-library/react";
import * as fc from "fast-check";
import { Markdown } from "./Markdown";

// ---------------------------------------------------------------------------
// Mocks
// ---------------------------------------------------------------------------

// Mock echarts-for-react to avoid canvas/DOM issues in jsdom
vi.mock("echarts-for-react", () => ({
  default: ({ option, className }: { option: any; className?: string }) => (
    <div data-testid="echarts-chart" data-option={JSON.stringify(option)} className={className}>
      ECharts Chart
    </div>
  ),
}));

// ---------------------------------------------------------------------------
// Arbitraries
// ---------------------------------------------------------------------------

/**
 * Generate JSON-serializable objects (not primitives).
 * fc.jsonValue() includes null, numbers, strings, booleans — but the Markdown
 * component only renders EChartsChart when JSON.parse returns a non-null value
 * AND the code block is detected as a block (match || isBlock). We restrict to
 * objects and arrays since those are the only valid ECharts option shapes.
 */
const jsonObjectArb = fc.oneof(
  fc.dictionary(fc.string(), fc.jsonValue()),
  fc.array(fc.jsonValue()),
);

/**
 * Generate strings that are NOT valid JSON.
 * Excludes empty strings since an empty echarts block produces no <code> element.
 * Includes malformed JSON and non-JSON text.
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
      return false; // Valid JSON, exclude it
    } catch {
      return true; // Invalid JSON, include it
    }
  }),
);

// ---------------------------------------------------------------------------
// Property-Based Tests
// ---------------------------------------------------------------------------

describe("Markdown ECharts rendering (Property-Based Tests)", () => {
  /**
   * Feature: echarts-markdown-renderer, Property 1:
   * Valid JSON echarts blocks render as charts in Markdown
   *
   * For any valid JSON object used as the body of an `echarts` fenced code block
   * rendered by `Markdown` in `minimal` or `full` mode, the output should contain
   * an `EChartsChart` component with that object as the `option` prop and
   * `className="echarts-chart-wrapper"`.
   *
   * Validates: Requirements 3.1, 3.2, 3.6
   */
  it("Property 1: for any fc.jsonValue() object, serialized and embedded in an echarts block, Markdown in minimal and full modes renders EChartsChart with the correct option and className='echarts-chart-wrapper'", () => {
    fc.assert(
      fc.property(
        jsonObjectArb,
        fc.constantFrom("minimal" as const, "full" as const),
        (jsonObj, mode) => {
          const jsonString = JSON.stringify(jsonObj);
          const markdown = `\`\`\`echarts\n${jsonString}\n\`\`\``;

          const { getByTestId } = render(<Markdown mode={mode}>{markdown}</Markdown>);

          // Assert EChartsChart is rendered
          const chart = getByTestId("echarts-chart");
          expect(chart).toBeInTheDocument();

          // Assert className is correct
          expect(chart).toHaveClass("echarts-chart-wrapper");

          // Assert option prop matches the original JSON object
          const optionAttr = chart.getAttribute("data-option");
          expect(optionAttr).not.toBeNull();
          const parsedOption = JSON.parse(optionAttr!);
          expect(parsedOption).toEqual(jsonObj);

          cleanup();
        },
      ),
      { numRuns: 100 },
    );
  });

  /**
   * Feature: echarts-markdown-renderer, Property 2:
   * Invalid JSON echarts blocks fall back to code in Markdown
   *
   * For any string that is not valid JSON used as the body of an `echarts` fenced
   * code block rendered by `Markdown`, the output should contain a fallback code
   * block with the original string content unchanged, and no `EChartsChart` component.
   *
   * Validates: Requirements 3.3
   */
  it("Property 2: for any non-JSON string (fc.string() filtered to exclude valid JSON), Markdown renders a fallback code block containing the original string and no EChartsChart", () => {
    fc.assert(
      fc.property(
        invalidJsonArb,
        fc.constantFrom("minimal" as const, "full" as const),
        (invalidJson, mode) => {
          const markdown = `\`\`\`echarts\n${invalidJson}\n\`\`\``;

          const { container, queryByTestId } = render(<Markdown mode={mode}>{markdown}</Markdown>);

          // Assert NO EChartsChart is rendered
          const chart = queryByTestId("echarts-chart");
          expect(chart).not.toBeInTheDocument();

          // Assert fallback code block is rendered and contains the original string.
          // react-markdown strips a trailing newline from code content, so we check
          // that the code element's text content matches the trimmed input.
          const codeBlock = container.querySelector("code");
          expect(codeBlock).toBeInTheDocument();
          expect(codeBlock?.textContent).toBe(invalidJson.replace(/\n$/, ""));

          cleanup();
        },
      ),
      { numRuns: 100 },
    );
  });

  /**
   * Additional test: terminal mode should NOT render ECharts
   *
   * Validates: Requirements 3.5
   */
  it("Property 2.1: terminal mode renders fallback code block for valid JSON echarts blocks", () => {
    fc.assert(
      fc.property(jsonObjectArb, (jsonObj) => {
        const jsonString = JSON.stringify(jsonObj);
        const markdown = `\`\`\`echarts\n${jsonString}\n\`\`\``;

        const { queryByTestId, container } = render(<Markdown mode="terminal">{markdown}</Markdown>);

        // Assert NO EChartsChart is rendered in terminal mode
        const chart = queryByTestId("echarts-chart");
        expect(chart).not.toBeInTheDocument();

        // Assert fallback code block is rendered
        const codeBlock = container.querySelector("code");
        expect(codeBlock).toBeInTheDocument();

        cleanup();
      }),
      { numRuns: 100 },
    );
  });
});
