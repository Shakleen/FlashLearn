import { useState } from "react";
import { useEffect } from "react";
import NavBar from "../components/NavBar";
import {
  useNavigate,
  useParams,
  useLocation,
  NavigateFunction,
} from "react-router-dom";
import { toast } from "sonner";
import Spinner from "../components/Spinner";
import { DeckItem } from "./DeckPage";

const FIELD_NAME = "deckName";
const FIELD_DESCRIPTION = "deckDescription";

function DeckFormPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [maxLengths, setMaxLengths] = useState<{
    name: number;
    description: number;
  }>({ name: -1, description: -1 });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchData(setLoading, setMaxLengths);
  }, []);

  const rightComponents = [
    <button type="submit" className="btn btn-light mx-2">
      Submit
    </button>,
    <button
      onClick={() => {
        navigate(-1);
      }}
      type="button"
      className="btn btn-primary mx-2"
    >
      Cancel
    </button>,
  ];

  return (
    <>
      <NavBar rightComponents={rightComponents} />
      {getBody(loading, setLoading, navigate, id, maxLengths)}
    </>
  );
}

function getBody(
  loading: boolean,
  setLoading: (loading: boolean) => void,
  navigate: NavigateFunction,
  id: string | undefined,
  maxLengths: { name: number; description: number }
) {
  if (loading) {
    return <Spinner />;
  }

  const location = useLocation();
  const deck = location.state?.deck as DeckItem;

  return (
    <form
      className="container m-3"
      onSubmit={(event) => handleSubmit(event, setLoading, navigate, id)}
    >
      <div className="mb-3">
        <label htmlFor="deckName" className="form-label">
          Deck Name
        </label>
        <input
          type="text"
          className="form-control"
          name={FIELD_NAME}
          aria-describedby="deckNameHelp"
          placeholder="Computer Science"
          defaultValue={deck?.name}
          maxLength={maxLengths.name}
          required
        />
        <div id="deckNameHelp" className="form-text">
          This is the name of the deck. <b>Max length: {maxLengths.name}</b>
        </div>
      </div>
      <div className="mb-3">
        <label htmlFor="deckDescription" className="form-label">
          Deck Description
        </label>
        <input
          type="text"
          className="form-control"
          name={FIELD_DESCRIPTION}
          placeholder="This a deck containing all the topics in computer science"
          defaultValue={deck?.description}
          maxLength={maxLengths.description}
        />
        <div id="deckDescriptionHelp" className="form-text">
          This is the description of the deck.{" "}
          <b>Max length: {maxLengths.description}</b>
        </div>
      </div>
    </form>
  );
}

async function HandlePostRequest(
  deckID: number,
  deckName: string,
  deckDescription: string
) {
  return fetch(
    "http://localhost:8080/deck" + (deckID == -1 ? "" : `/${deckID}`),
    {
      method: "POST",
      body: JSON.stringify({
        name: deckName,
        description: deckDescription,
      }),
    }
  );
}

async function fetchData(
  setLoading: (loading: boolean) => void,
  setMaxLengths: (maxLengths: { name: number; description: number }) => void
) {
  setLoading(true);

  try {
    const responseNameMaxLength = await fetch(
      "http://localhost:8080/deck/nameMaxLength"
    );
    const resultNameMaxLength = await responseNameMaxLength.json();

    const responseDescriptionMaxLength = await fetch(
      "http://localhost:8080/deck/descriptionMaxLength"
    );
    const resultDescriptionMaxLength =
      await responseDescriptionMaxLength.json();

    setMaxLengths({
      name: resultNameMaxLength["maxLength"],
      description: resultDescriptionMaxLength["maxLength"],
    });
  } catch (error) {
    toast.error(
      `Error fetching data: ${
        error instanceof Error ? error.message : "Unknown error"
      }`
    );
  } finally {
    setLoading(false);
  }
}

async function handleSubmit(
  event: React.FormEvent<HTMLFormElement>,
  setLoading: (loading: boolean) => void,
  navigate: NavigateFunction,
  id: string | undefined
) {
  event.preventDefault();
  setLoading(true);

  try {
    const formData = new FormData(event.target as HTMLFormElement);
    const payload = Object.fromEntries(formData);
    const deckName = payload[FIELD_NAME];
    const deckDescription = payload[FIELD_DESCRIPTION];

    const deckID: number = parseInt(id || "-1");

    const response = await HandlePostRequest(
      deckID,
      deckName as string,
      deckDescription as string
    );

    if (!response.ok) {
      if (response.status === 409) {
        throw new Error("Deck name already exists");
      } else {
        throw new Error(`HTTP Error`);
      }
    }

    const toastMessage: string = `Deck ${
      deckID == -1 ? "created" : "updated"
    } successfully`;
    toast.success(toastMessage);

    navigate(-1);
  } catch (error) {
    toast.error(
      `Error submitting form: ${
        error instanceof Error ? error.message : "Unknown error"
      }`
    );
  } finally {
    setLoading(false);
  }
}

export default DeckFormPage;
