import NavBar from "../components/NavBar";
import { useNavigate } from "react-router-dom";

function DeckFormPage() {
  const navigate = useNavigate();

  return (
    <>
      <NavBar />
      <form className="container m-3">
        <div className="mb-3">
          <label htmlFor="deckName" className="form-label">
            Deck Name
          </label>
          <input
            type="text"
            className="form-control"
            id="deckName"
            aria-describedby="deckNameHelp"
            placeholder="Computer Science"
          />
          <div id="deckNameHelp" className="form-text">
            This is the name of the deck.
          </div>
        </div>
        <div className="mb-3">
          <label htmlFor="deckDescription" className="form-label">
            Deck Description
          </label>
          <input
            type="text"
            className="form-control"
            id="deckDescription"
            placeholder="This a deck containing all the topics in computer science"
          />
          <div id="deckDescriptionHelp" className="form-text">
            This is the description of the deck.
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
