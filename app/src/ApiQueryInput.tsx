import React, { useState } from 'react';

interface ApiQueryInputProps {
  onSubmit: (query: string) => void;
}

const ApiQueryInput: React.FC<ApiQueryInputProps> = ({ onSubmit }) => {
  const [query, setQuery] = useState('');

  const handleQueryChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setQuery(event.target.value);
  };

  const handleSubmit = () => {
    onSubmit(query);
  };

  return (
    <div className="flex items-center justify-center h-64 border border-gray-300 rounded-lg p-4">
      <label htmlFor="apiQuery">API Query: </label>
      <input
        type="text"
        id="apiQuery"
        value={query}
        onChange={handleQueryChange}
        placeholder="Enter your API query here"
        className="dark-input"
      />
      <button type="button" onClick={handleSubmit}>Submit</button>
    </div>
  );
};

export default ApiQueryInput;
