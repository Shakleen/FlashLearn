import { useState } from "react";
import { useEffect } from "react";
import NavBar from "../components/NavBar";
import { useNavigate, useParams } from "react-router-dom";
import { toast } from "sonner";

function DeckFormPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [maxLengths, setMaxLengths] = useState<{
    name: number;
    description: number;
  }>({ name: -1, description: -1 });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
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
    };

    fetchData();
  }, []);

  if (loading) {
    return <div>Loading...</div>;
  }

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setLoading(true);

    try {
      const formData = new FormData(event.target as HTMLFormElement);
      const payload = Object.fromEntries(formData);
      const deckName = payload["deckName"];
      const deckDescription = payload["deckDescription"];

      const deckID: number = parseInt(id || "-1");

      const response = await HandlePostRequest(
        deckID,
        deckName as string,
        deckDescription as string
      );

      if (!response.ok) {
        throw new Error(`Http error! Status: ${response.status}`);
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
  };

  return (
    <>
      <NavBar />
      <form className="container m-3" onSubmit={handleSubmit}>
        <div className="mb-3">
          <label htmlFor="deckName" className="form-label">
            Deck Name
          </label>
          <input
            type="text"
            className="form-control"
            name="deckName"
            aria-describedby="deckNameHelp"
            placeholder="Computer Science"
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
            name="deckDescription"
            placeholder="This a deck containing all the topics in computer science"
            maxLength={maxLengths.description}
          />
          <div id="deckDescriptionHelp" className="form-text">
            This is the description of the deck.{" "}
            <b>Max length: {maxLengths.description}</b>
          </div>
        </div>
        <button type="submit" className="btn btn-primary">
          Submit
        </button>
        <button
          onClick={() => {
            navigate(-1);
          }}
          type="button"
          className="btn btn-secondary"
        >
          Cancel
        </button>
      </form>
    </>
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

export default DeckFormPage;
