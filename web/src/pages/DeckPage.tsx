import { useParams } from "react-router-dom";
import { useEffect, useState } from "react";

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

  return (
    <>
      <nav className="navbar bg-body-tertiary" data-bs-theme="dark">
        <div className="container-fluid">
          <a className="navbar-brand" href="#">
            Flash Learn
          </a>
        </div>
      </nav>
      <h1>Deck: {result.name}</h1>
      <p>Description: {result.description}</p>
      <p>Created: {result.creationDate}</p>
      <p>Modified: {result.modificationDate}</p>
      <p>Last Studied: {result.lastStudyDate}</p>
      <p>Total Cards: {result.totalCards}</p>
    </>
  );
}

export default DeckPage;
