import React, { useEffect, useRef, useState } from 'react';
import Chart from 'chart.js/auto';
import zoomPlugin from 'chartjs-plugin-zoom';

interface FloatParams {
  [key: string]: number[];
}

interface IntParams {
  [key: string]: number[];
}

interface JsonLogEntry {
  partition_index: number;
  objective: number;
  float_params: FloatParams;
  int_params: IntParams;
}

interface LineChartProps {
  data: {
    log_filename: string;
    partition_iterations: number;
    entry: JsonLogEntry;
  }[];
}

const LineChart: React.FC<LineChartProps> = ({ data }) => {
  const chartRef = useRef<HTMLCanvasElement | null>(null);
  const chartInstanceRef = useRef<Chart | null>(null);
  const [lineColours, setLineColours] = useState<{
    [filenamePartitionIndex: string]: string
  }>({});
  Chart.register(zoomPlugin);

  function getRandomColour() {
    const letters = '0123456789ABCDEF';
    let color = '#';
    const minBrightness = 0.6;
  
    for (let i = 0; i < 6; i++) {
      color += letters[Math.floor(Math.random() * 16)];
    }
  
    // Get the RGB components of the color
    const red = parseInt(color.substring(1, 3), 16);
    const green = parseInt(color.substring(3, 5), 16);
    const blue = parseInt(color.substring(5, 7), 16);
  
    // Calculate the brightness of the color (normalized value between 0 and 1)
    const brightness = (0.299 * red + 0.587 * green + 0.114 * blue) / 255;
  
    // If the brightness is below the minimum, adjust the color to make it brighter
    if (brightness < minBrightness) {
      const correctionFactor = minBrightness / brightness;
      const newRed = Math.min(255, Math.floor(red * correctionFactor));
      const newGreen = Math.min(255, Math.floor(green * correctionFactor));
      const newBlue = Math.min(255, Math.floor(blue * correctionFactor));
  
      // Convert the RGB components back to hexadecimal and update the color
      color = `#${(newRed < 16 ? '0' : '') + newRed.toString(16)}` +
              `${(newGreen < 16 ? '0' : '') + newGreen.toString(16)}` +
              `${(newBlue < 16 ? '0' : '') + newBlue.toString(16)}`;
    }
  
    return color;
  }

  const processChartData = (data: LineChartProps['data']) => {
    const result: { [filenamePartitionIndex: string]: {
      label: string;
      data: {
        x: number;
        y: number;
      }[];
      borderColor: string;
      borderWidth: number;
      fill: boolean;
    }} = {};

    if (!data || data.length === 0) {
      return {}; // Return an empty object if there's no data
    }

    data.forEach((item) => { 
      const filenamePartitionIndex = item.log_filename + 
        " " + String(item.entry.partition_index);
      const objectiveData = item.entry.objective;

      if (!(filenamePartitionIndex in lineColours)) {
        const colour = getRandomColour()
          setLineColours((prevLineColours) => ({
            ...prevLineColours,
            [filenamePartitionIndex]: colour,
          }));
      }

      if (!(filenamePartitionIndex in result)) {
        result[filenamePartitionIndex] = {
          label: `${filenamePartitionIndex}`,
          data: [],
          borderColor: lineColours[filenamePartitionIndex],
          borderWidth: 2,
          fill: false,
        };
      }

      result[filenamePartitionIndex].data.push({
        x: item.partition_iterations,
        y: objectiveData,
      })
    });

    return result;
  };

  const chartData = {
    datasets: Object.values(processChartData(data)).flat()
  };

  const resetZoom = () => {
    if (chartInstanceRef.current) {
      chartInstanceRef.current.resetZoom();
    }
  };

  useEffect(() => {
    const handleKeyPress = (event: KeyboardEvent) => {
      if (event.key === 'r') {
        resetZoom();
      }
    };
  
    window.addEventListener('keypress', handleKeyPress);

    return () => {
      window.removeEventListener('keypress', handleKeyPress);
    };
  }, [resetZoom]);

  useEffect(() => {
    if (!chartRef.current) return;

    if (chartInstanceRef.current) {
      chartInstanceRef.current.destroy();
    }

    const ctx = chartRef.current.getContext('2d');
    if (ctx) {
      chartRef.current.height = 300;
      chartInstanceRef.current = new Chart(ctx, {
        type: 'line',
        data: chartData,
        options: {
          responsive: true,
          maintainAspectRatio: false,
          animation: {
            duration: 0
          },
          scales: {
            x: {
              type: 'linear',
              position: 'bottom',
              grid: {
                display: false
              },
              ticks: {
                color: 'white'
              },
            },
            y: {
              display: true,
              grid: {
                display: false
              },
              ticks: {
                color: 'white'
              },
            },
          },
          plugins: {
            legend: {
              labels: {
                color: 'white'
              }
            },
            title: {
              display: false
            },
            zoom: {
              pan: {
                enabled: true,
                mode: 'xy',
                modifierKey: 'ctrl',
              },
              zoom: {
                drag: {
                  enabled: true,
                  modifierKey: 'shift',
                },
                wheel: {
                  enabled: true,
                },
                pinch: {
                  enabled: true
                },
                mode: 'xy',
              },
            },
          },
          elements: {
            point: {
              borderColor: 'white',
              borderWidth: 1
            },
            line: {
              borderColor: 'white',
              borderWidth: 1
            }
          }
        },
      });
    }
  }, [chartData]);

  return (
    <div className="flex items-center justify-center h-64 border border-gray-300 rounded-lg p-4">
      <canvas ref={chartRef} width="400" height="200" />
    </div>
  );
};

export default LineChart;
