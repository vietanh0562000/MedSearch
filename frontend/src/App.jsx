import { useState } from "react";
import { meds } from "./meds";

function App() {
  const [query, setQuery] = useState("");
  const filteredMeds = meds.filter(med =>
    med.name.toLowerCase().includes(query.toLowerCase())
  );

  return (
    <div style={{
      minHeight: "100vh",
      display: "flex",
      flexDirection: "column",
      justifyContent: "center",
      alignItems: "center",
      fontFamily: "sans-serif"
    }}>
      <div style={{ width: 500, maxWidth: "90%" }}>
        <h1 style={{ textAlign: "center" }}>Med Search</h1>
        <input
          type="text"
          placeholder="Search medication..."
          value={query}
          onChange={e => setQuery(e.target.value)}
          style={{
            width: "100%",
            padding: "0.5rem",
            fontSize: "1rem",
            marginBottom: "1rem"
          }}
        />
        <ul style={{ padding: 0, listStyle: "none" }}>
          {filteredMeds.length === 0 && <li style={{ textAlign: "center" }}>No results found.</li>}
          {filteredMeds.map(med => (
            <li key={med.name} style={{ marginBottom: "1rem", textAlign: "center" }}>
              <strong>{med.name}</strong>
              <div>{med.description}</div>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}

export default App;
