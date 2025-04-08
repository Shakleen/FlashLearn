import { Link, useParams } from "react-router-dom";
import { useEffect, useState } from "react";
import NavBar from "../components/NavBar";

export interface DeckItem {
  id: string;
  name: string;
  description: string;
  creationDate: string;
  modificationDate: string;
  lastStudyDate: string;
  totalCards: number;
}

function DeckPage() {
  const { id } = useParams();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [result, setResult] = useState<DeckItem | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      setError(null);
      try {
        const response = await fetch(`http://localhost:8080/deck/${id}`);

        if (!response.ok) {
          throw new Error(`Http error! Status: ${response.status}`);
        }

        const rawData = await response.json();

        const deckData: DeckItem = {
          id: rawData.id || "",
          name: rawData.name || "",
          description: rawData.description || "",
          creationDate: rawData.creation_date || "",
          modificationDate: rawData.modification_date || "",
          lastStudyDate: rawData.last_study_date || "",
          totalCards: rawData.total_cards || 0,
        };

        setResult(deckData);
      } catch (error) {
        setError(
          `Error fetching data: ${
            error instanceof Error ? error.message : "Unknown error"
          }`
        );
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [id]);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>{error}</div>;
  }

  if (!result) {
    return <div>No deck data found.</div>;
  }

  const creationDateString = new Date(result.creationDate).toLocaleDateString();
  const modificationDateString = new Date(
    result.modificationDate
  ).toLocaleDateString();

  let lastStudyDateString = "Never";
  const lastStudyYear = new Date(result.lastStudyDate).getFullYear();
  if (lastStudyYear > 2000) {
    lastStudyDateString = new Date(result.lastStudyDate).toLocaleDateString();
  }

  return (
    <>
      <NavBar />
      <div className="card text-center m-2">
        <div className="card-header">Last Studied: {lastStudyDateString}</div>
        <div className="card-body">
          <h5 className="card-title">{result.name}</h5>
          <p className="card-text">{result.description}</p>
          <p className="card-text">
            <small className="text-body-secondary">
              Created: {creationDateString} <br />
              Last updated: {modificationDateString} <br />
              Total cards: {result.totalCards}
            </small>
          </p>
          <div
            className="btn-group mb-3"
            role="group"
            aria-label="Basic mixed styles example"
          >
            <Link to={"/"} className="btn btn-primary">
              Study Now
            </Link>
            <button type="button" className="btn btn-secondary">
              Edit
            </button>
            <button type="button" className="btn btn-danger">
              Delete
            </button>
          </div>
        </div>
      </div>
    </>
  );
}

export default DeckPage;
