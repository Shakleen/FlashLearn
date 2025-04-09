import { useState } from "react";
import { useEffect } from "react";
import NavBar from "../components/NavBar";
import { useNavigate, useParams } from "react-router-dom";

function DeckFormPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [nameMaxLength, setNameMaxLength] = useState(0);
  const [descriptionMaxLength, setDescriptionMaxLength] = useState(0);

  useEffect(() => {
    const fetchData = async () => {
      const responseNameMaxLength = await fetch(
        "http://localhost:8080/deck/nameMaxLength"
      );
      const resultNameMaxLength = await responseNameMaxLength.json();
      setNameMaxLength(resultNameMaxLength["maxLength"]);

      const responseDescriptionMaxLength = await fetch(
        "http://localhost:8080/deck/descriptionMaxLength"
      );
      const resultDescriptionMaxLength =
        await responseDescriptionMaxLength.json();
      setDescriptionMaxLength(resultDescriptionMaxLength["maxLength"]);
    };

    fetchData();
  }, []);

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const formData = new FormData(event.target as HTMLFormElement);
    const payload = Object.fromEntries(formData);
    const deckName = payload["deckName"];
    const deckDescription = payload["deckDescription"];

    var deckID = parseInt(id || "-1");
    console.log(deckID);

    if (deckID == -1) {
      // Insert POST /deck
      const response = await fetch("http://localhost:8080/deck", {
        method: "POST",
        body: JSON.stringify({
          name: deckName,
          description: deckDescription,
        }),
      });

      if (!response.ok) {
        throw new Error(`Http error! Status: ${response.status}`);
      }

      const rawData = await response.json();
      deckID = rawData["id"];
    } else {
      // Update POST /deck/{id}
      const response = await fetch(`http://localhost:8080/deck/${deckID}`, {
        method: "POST",
        body: JSON.stringify({
          name: deckName,
          description: deckDescription,
        }),
      });

      if (!response.ok) {
        throw new Error(`Http error! Status: ${response.status}`);
      }

      const rawData = await response.json();
      deckID = rawData["id"];
    }

    navigate(-1);
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
            maxLength={nameMaxLength}
            required
          />
          <div id="deckNameHelp" className="form-text">
            This is the name of the deck. <b>Max length: {nameMaxLength}</b>
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
            maxLength={descriptionMaxLength}
          />
          <div id="deckDescriptionHelp" className="form-text">
            This is the description of the deck.{" "}
            <b>Max length: {descriptionMaxLength}</b>
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

export default DeckFormPage;
