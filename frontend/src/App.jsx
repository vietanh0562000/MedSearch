import { useState } from "react";

function App() {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSearch = async () => {
    if (!query.trim()) return;
    
    setLoading(true);
    setError("");
    
    try {
      const response = await fetch(`http://localhost:8080/v1/api/search?text=${encodeURIComponent(query)}`);
      console.log('Response status:', response.status);
      console.log('Response ok:', response.ok);
      console.log('Response headers:', response.headers);
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();
      console.log(data)
      setResults(data);
    } catch (err) {
      setError("Failed to fetch results. Please try again.");
      console.error("Search error:", err);
    } finally {
      setLoading(false);
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === "Enter") {
      handleSearch();
    }
  };

  return (
    <div className="app-container">
      <div className="app-content">
        <h1 className="app-title">Med Search</h1>
        <div className="search-container">
          <input
            type="text"
            placeholder="Search medication..."
            value={query}
            onChange={e => setQuery(e.target.value)}
            onKeyPress={handleKeyPress}
            className="search-input"
          />
          <button
            onClick={handleSearch}
            disabled={loading}
            className={`search-button ${loading ? 'loading' : ''}`}
          >
            {loading ? "Searching..." : "Search"}
          </button>
        </div>
        
        {error && (
          <div className="error-message">
            {error}
          </div>
        )}
        
        <ul className="results-list">
          {results.length === 0 && !loading && !error && (
            <li className="no-results">
              Enter a search term and press Enter to search
            </li>
          )}
          {results.map((med, index) => (
            <li key={index} className="result-item">
              <strong>{med.Name}</strong>
              {med.description && (
                <div className="result-description">
                  {med.Description}
                </div>
              )}
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}

export default App;
