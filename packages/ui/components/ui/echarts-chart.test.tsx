import { describe, it, expect, vi } from "vitest";
import { render } from "@testing-library/react";
import type { EChartsOption } from "echarts";
import { EChartsChart } from "./echarts-chart";
import { Markdown } from "../../markdown/Markdown";

// ---------------------------------------------------------------------------
// Mocks
// ---------------------------------------------------------------------------

// Mock echarts-for-react to avoid canvas/DOM issues in jsdom
vi.mock("echarts-for-react", () => ({
  default: ({ option, className, style, theme, notMerge, lazyUpdate, showLoading }: {
    option: unknown;
    className?: string;
    style?: React.CSSProperties;
    theme?: string | Record<string, unknown>;
    notMerge?: boolean;
    lazyUpdate?: boolean;
    showLoading?: boolean;
  }) => (
    <div
      data-testid="echarts-chart"
      data-option={JSON.stringify(option)}
      data-theme={typeof theme === "string" ? theme : JSON.stringify(theme)}
      data-not-merge={String(notMerge)}
      data-lazy-update={String(lazyUpdate)}
      data-show-loading={String(showLoading)}
      className={className}
      style={style}
    >
      ECharts Chart
    </div>
  ),
}));

// ---------------------------------------------------------------------------
// EChartsChart component unit tests
// ---------------------------------------------------------------------------

describe("EChartsChart", () => {
  /**
   * 8.2 — Renders without errors given a minimal valid option
   */
  it("renders without errors given a minimal valid option", () => {
    const option: EChartsOption = { series: [{ type: "bar", data: [1, 2, 3] }] };
    const { getByTestId } = render(<EChartsChart option={option} />);
    const chart = getByTestId("echarts-chart");
    expect(chart).toBeInTheDocument();
  });

  it("passes the option prop to the underlying chart", () => {
    const option: EChartsOption = { series: [{ type: "bar", data: [1, 2, 3] }] };
    const { getByTestId } = render(<EChartsChart option={option} />);
    const chart = getByTestId("echarts-chart");
    expect(JSON.parse(chart.getAttribute("data-option")!)).toEqual(option);
  });

  /**
   * 8.3 — Accepts all optional props without throwing
   */
  it("accepts all optional props without throwing", () => {
    const option: EChartsOption = { series: [{ type: "line", data: [10, 20, 30] }] };
    const onChartReady = vi.fn();
    const onEvents = { click: vi.fn() };

    expect(() =>
      render(
        <EChartsChart
          option={option}
          className="my-chart"
          style={{ height: 300 }}
          theme="dark"
          onEvents={onEvents}
          notMerge={true}
          lazyUpdate={false}
          showLoading={false}
          onChartReady={onChartReady}
        />
      )
    ).not.toThrow();
  });

  it("forwards className to the underlying chart element", () => {
    const option = { series: [] };
    const { getByTestId } = render(<EChartsChart option={option} className="my-custom-class" />);
    expect(getByTestId("echarts-chart")).toHaveClass("my-custom-class");
  });

  it("forwards notMerge, lazyUpdate, and showLoading props", () => {
    const option = { series: [] };
    const { getByTestId } = render(
      <EChartsChart option={option} notMerge={true} lazyUpdate={true} showLoading={true} />
    );
    const chart = getByTestId("echarts-chart");
    expect(chart.getAttribute("data-not-merge")).toBe("true");
    expect(chart.getAttribute("data-lazy-update")).toBe("true");
    expect(chart.getAttribute("data-show-loading")).toBe("true");
  });

  it("renders with a complex option (xAxis, yAxis, series)", () => {
    const option: EChartsOption = {
      xAxis: { type: "category", data: ["Mon", "Tue", "Wed"] },
      yAxis: { type: "value" },
      series: [{ type: "bar", data: [120, 200, 150] }],
    };
    const { getByTestId } = render(<EChartsChart option={option} />);
    expect(getByTestId("echarts-chart")).toBeInTheDocument();
  });
});

// ---------------------------------------------------------------------------
// Markdown integration tests (8.4 and 8.5)
// ---------------------------------------------------------------------------

describe("Markdown ECharts integration", () => {
  const validOption: EChartsOption = { series: [{ type: "bar", data: [1, 2, 3] }] };
  const validEchartsBlock = `\`\`\`echarts\n${JSON.stringify(validOption)}\n\`\`\``;

  /**
   * 8.4 — terminal mode renders a plain code block, NOT EChartsChart
   */
  it("mode='terminal' renders a plain code block for an echarts block", () => {
    const { queryByTestId, container } = render(
      <Markdown mode="terminal">{validEchartsBlock}</Markdown>
    );
    // No EChartsChart in terminal mode
    expect(queryByTestId("echarts-chart")).not.toBeInTheDocument();
    // A code element should be present as fallback
    expect(container.querySelector("code")).toBeInTheDocument();
  });

  /**
   * 8.5 — minimal mode renders EChartsChart for a valid echarts block
   */
  it("mode='minimal' renders EChartsChart for a valid echarts block", () => {
    const { getByTestId } = render(
      <Markdown mode="minimal">{validEchartsBlock}</Markdown>
    );
    const chart = getByTestId("echarts-chart");
    expect(chart).toBeInTheDocument();
    expect(chart).toHaveClass("echarts-chart-wrapper");
    expect(JSON.parse(chart.getAttribute("data-option")!)).toEqual(validOption);
  });

  /**
   * 8.5 — full mode renders EChartsChart for a valid echarts block
   */
  it("mode='full' renders EChartsChart for a valid echarts block", () => {
    const { getByTestId } = render(
      <Markdown mode="full">{validEchartsBlock}</Markdown>
    );
    const chart = getByTestId("echarts-chart");
    expect(chart).toBeInTheDocument();
    expect(chart).toHaveClass("echarts-chart-wrapper");
    expect(JSON.parse(chart.getAttribute("data-option")!)).toEqual(validOption);
  });

  it("mode='minimal' renders fallback code block for invalid JSON echarts block", () => {
    const invalidBlock = "```echarts\n{not valid json\n```";
    const { queryByTestId, container } = render(
      <Markdown mode="minimal">{invalidBlock}</Markdown>
    );
    expect(queryByTestId("echarts-chart")).not.toBeInTheDocument();
    expect(container.querySelector("code")).toBeInTheDocument();
  });

  it("mode='full' renders fallback code block for invalid JSON echarts block", () => {
    const invalidBlock = "```echarts\n{not valid json\n```";
    const { queryByTestId, container } = render(
      <Markdown mode="full">{invalidBlock}</Markdown>
    );
    expect(queryByTestId("echarts-chart")).not.toBeInTheDocument();
    expect(container.querySelector("code")).toBeInTheDocument();
  });
});
