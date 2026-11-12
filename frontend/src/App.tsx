import { Calculator } from "./components/Calculator";
import "./App.css";

export function App() {
  return (
    <div className="page">
      <header className="page-header">
        <p className="eyebrow">Sezzle</p>
        <h1>Sezzle Calculator</h1>
        <p className="lede">
          Add, subtract, multiply, and divide. Every result is computed by the API.
        </p>
      </header>
      <main>
        <Calculator />
      </main>
    </div>
  );
}
