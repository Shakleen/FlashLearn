import { Link, useParams } from "react-router-dom";
import { useEffect, useState } from "react";
import NavBar from "../components/NavBar";
import Spinner from "../components/Spinner";

export interface DeckItem {
  id: string;
  name: string;
  description: string;
  creationDate: string;
  modificationDate: string;
  lastStudyDate: string;
}

function DeckPage() {
  const { id } = useParams();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [result, setResult] = useState<DeckItem | null>(null);

  useEffect(() => {
    fetchData(setLoading, setError, setResult, id);
  }, [id]);

  const rightComponents = [
    <Link
      to={`/card/form/${id}`}
      className="btn btn-light"
      state={{ deck: result }}
    >
      Add Cards
    </Link>,
    <Link
      to={`/deck/form/${id}`}
      state={{ deck: result }}
      className="btn btn-light"
    >
      Edit
    </Link>,
    <Link to={`/deck/delete/${id}`} className="btn btn-danger">
      Delete
    </Link>,
  ];

  return (
    <>
      <NavBar rightComponents={rightComponents} />
      {getBody(loading, error, result)}
    </>
  );
}

async function fetchData(
  setLoading: (loading: boolean) => void,
  setError: (error: string | null) => void,
  setResult: (result: DeckItem | null) => void,
  id: string | undefined
) {
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
}

function getBody(
  loading: boolean,
  error: string | null,
  result: DeckItem | null
) {
  var body;

  if (loading) {
    body = <Spinner />;
  } else if (error) {
    body = <div>{error}</div>;
  } else if (!result) {
    body = <div>No deck data found.</div>;
  } else {
    const creationDateString = new Date(
      result.creationDate
    ).toLocaleDateString();
    const modificationDateString = new Date(
      result.modificationDate
    ).toLocaleDateString();

    let lastStudyDateString = "Never";
    const lastStudyYear = new Date(result.lastStudyDate).getFullYear();
    if (lastStudyYear > 2000) {
      lastStudyDateString = new Date(result.lastStudyDate).toLocaleDateString();
    }

    body = (
      <div className="card text-center m-2">
        <div className="card-header">Last Studied: {lastStudyDateString}</div>
        <div className="card-body">
          <h5 className="card-title fw-bold fs-2">{result.name}</h5>
          <p className="card-text fs-5">{result.description}</p>
          <p className="card-text">
            <small className="text-body-secondary">
              Created: {creationDateString} <br />
              Last updated: {modificationDateString}
            </small>
          </p>
          <Link
            to={"/"}
            className="btn btn-primary mx-2"
            style={{ width: "120px" }}
          >
            Study Now
          </Link>
        </div>
      </div>
    );
  }
  return body;
}

export default DeckPage;
