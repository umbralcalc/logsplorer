import React, { useState, useEffect } from 'react';
import ApiQueryInput from './ApiQueryInput';
import LineChart from './LineChart';

const App: React.FC = () => {
  const [apiQuery, setApiQuery] = useState('');
  const [chartData, setChartData] = useState<any>(null);

  const handleApiQuerySubmit = (query: string) => {
    setApiQuery(query);
  };

  const fetchData = async () => {
    try {
      const response = await fetch(`http://localhost:8080/api/logsplorer?${apiQuery}`);
      const responseData = await response.json();
      console.log(responseData)
      setChartData(responseData);
    } catch (error) {
      console.error('Error fetching data:', error);
    }
  };

  useEffect(() => {
    fetchData();
  }, [apiQuery]);

  return (
    <div>
      <ApiQueryInput onSubmit={handleApiQuerySubmit} />
      <LineChart data={chartData} />
      <h5>Zoom: Mousewheel or Shift + Click + Drag</h5>
      <h5>Pan: Ctrl + Click + Drag</h5>
      <h5>Reset: r Key</h5>
    </div>
  );
};

export default App;
