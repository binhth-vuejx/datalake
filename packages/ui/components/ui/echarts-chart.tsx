"use client"

import type { EChartsOption } from 'echarts'
import type { EChartsReactProps } from 'echarts-for-react'
import ReactECharts from 'echarts-for-react'

export interface EChartsChartProps {
  /** Required ECharts option object */
  option: EChartsOption
  /** Additional CSS class names */
  className?: string
  /** Inline styles */
  style?: React.CSSProperties
  /** ECharts theme name or object */
  theme?: string | Record<string, unknown>
  /** Event handlers map */
  onEvents?: EChartsReactProps['onEvents']
  /** Whether to not merge with previous option (default: false) */
  notMerge?: boolean
  /** Whether to update lazily (default: false) */
  lazyUpdate?: boolean
  /** Whether to show loading animation */
  showLoading?: boolean
  /** Callback when chart instance is ready */
  onChartReady?: EChartsReactProps['onChartReady']
}

/**
 * Default option merged into every chart unless the caller explicitly overrides.
 * - tooltip: show on hover (trigger: 'axis' for line/bar, 'item' for pie/scatter)
 * - legend: auto-show when series names are present
 * - toolbox: save-as-image + data-view + zoom buttons
 * - grid: sensible padding so axis labels are never clipped
 *
 * The caller's option is spread on top, so any key they provide wins.
 */
const DEFAULT_OPTION: Partial<EChartsOption> = {
  tooltip: {
    trigger: 'axis',
    axisPointer: { type: 'cross' },
  },
  legend: {
    type: 'scroll',
    bottom: 0,
  },
  toolbox: {
    right: 8,
    feature: {
      dataZoom: { yAxisIndex: 'none' },
      restore: {},
      saveAsImage: {},
    },
  },
  grid: {
    left: '3%',
    right: '4%',
    bottom: 40,
    containLabel: true,
  },
}

export function EChartsChart({
  option,
  className,
  style,
  theme,
  onEvents,
  notMerge,
  lazyUpdate,
  showLoading,
  onChartReady,
}: EChartsChartProps) {
  // Merge defaults first, then caller option on top so explicit keys always win.
  // Deep-merge is intentionally avoided — if the caller provides `tooltip` they
  // own it entirely; we only fill in keys they omitted.
  const mergedOption: EChartsOption = { ...DEFAULT_OPTION, ...option }

  return (
    <ReactECharts
      option={mergedOption}
      className={className}
      style={{ width: '100%', height: 360, ...style }}
      theme={theme}
      onEvents={onEvents}
      notMerge={notMerge ?? true}
      lazyUpdate={lazyUpdate}
      showLoading={showLoading}
      onChartReady={onChartReady}
    />
  )
}
