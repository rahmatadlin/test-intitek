<template>
  <div class="chart-container">
    <canvas ref="chartCanvas"></canvas>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from "vue";
import {
  Chart,
  BarController,
  BarElement,
  CategoryScale,
  LinearScale,
  Title,
  Tooltip,
  Legend,
} from "chart.js";

// Register Chart.js components
Chart.register(
  BarController,
  BarElement,
  CategoryScale,
  LinearScale,
  Title,
  Tooltip,
  Legend
);

const props = defineProps({
  data: {
    type: Object,
    required: true,
  },
  options: {
    type: Object,
    default: () => ({}),
  },
  horizontal: {
    type: Boolean,
    default: false,
  },
});

const chartCanvas = ref(null);
let chartInstance = null;

const createChart = () => {
  if (!chartCanvas.value) return;

  // Destroy existing chart if it exists
  if (chartInstance) {
    chartInstance.destroy();
  }

  const ctx = chartCanvas.value.getContext("2d");

  const baseOptions = {
    responsive: true,
    maintainAspectRatio: false,
    indexAxis: props.horizontal ? "y" : "x",
    plugins: {
      legend: {
        display: false,
      },
      tooltip: {
        backgroundColor: "rgba(0, 0, 0, 0.8)",
        padding: 12,
        titleFont: {
          size: 14,
        },
        bodyFont: {
          size: 13,
        },
        callbacks: {
          label: function (context) {
            const value = props.horizontal
              ? context.parsed.x
              : context.parsed.y;
            return `Quantity: ${value}`;
          },
        },
      },
    },
    scales: props.horizontal
      ? {
          x: {
            beginAtZero: true,
            ticks: {
              precision: 0,
            },
          },
          y: {
            grid: {
              display: false,
            },
          },
        }
      : {
          x: {
            grid: {
              display: false,
            },
          },
          y: {
            beginAtZero: true,
            ticks: {
              precision: 0,
            },
          },
        },
    ...props.options,
  };

  chartInstance = new Chart(ctx, {
    type: "bar",
    data: props.data,
    options: baseOptions,
  });
};

onMounted(() => {
  createChart();
});

watch(
  () => props.data,
  () => {
    createChart();
  },
  { deep: true }
);
</script>

<style scoped>
.chart-container {
  position: relative;
  height: 100%;
  width: 100%;
}
</style>
