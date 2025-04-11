import { useLocation, useNavigate, useParams } from "react-router-dom";
import NavBar from "../components/NavBar";
import { toast } from "sonner";
import { DeckItem } from "./DeckPage";

const FIELD_FRONT = "front";
const FIELD_BACK = "back";
const FIELD_SOURCE = "source";
function CardFormPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const location = useLocation();
  const deck = location.state?.deck as DeckItem;
  const rightComponents = [
    <button
      onClick={() => {
        navigate(-1);
      }}
      type="button"
      className="btn btn-primary mx-2"
    >
      Go Back
    </button>,
  ];

  return (
    <>
      <NavBar rightComponents={rightComponents} />
      {getBody(id, deck)}
    </>
  );
}

function getBody(deckID: string | undefined, deck: DeckItem) {
  return (
    <>
      <div className="text-center fs-2 my-3">
        {deck.name}
        <p className="fs-5">{deck.description}</p>
      </div>
      <form
        className="container m-3"
        onSubmit={(event) => handleSubmit(event, deckID)}
      >
        <div className="mb-3">
          <label htmlFor="front" className="form-label">
            Front
          </label>
          <input
            type="text"
            className="form-control"
            name={FIELD_FRONT}
            aria-describedby="frontHelp"
            placeholder="What is the capital of France?"
            required
          />
          <div id="frontHelp" className="form-text">
            This is what you will see on the front of the card.
          </div>
        </div>
        <div className="mb-3">
          <label htmlFor="back" className="form-label">
            Back
          </label>
          <input
            type="text"
            className="form-control"
            name={FIELD_BACK}
            placeholder="Paris"
            required
          />
          <div id="backHelp" className="form-text">
            This is what you will see on the back of the card.
          </div>
        </div>
        <div className="mb-3">
          <label htmlFor="source" className="form-label">
            Source
          </label>
          <input
            type="text"
            className="form-control"
            name={FIELD_SOURCE}
            placeholder="http://source.com"
          />
          <div id="sourceHelp" className="form-text">
            This is where you found the information.
          </div>
        </div>
        <button type="submit" className="btn btn-primary">
          Submit
        </button>
      </form>
    </>
  );
}

async function handleSubmit(
  event: React.FormEvent<HTMLFormElement>,
  deckID: string | undefined
) {
  event.preventDefault();

  try {
    const formData = new FormData(event.target as HTMLFormElement);
    const payload = Object.fromEntries(formData);
    const front = payload[FIELD_FRONT] as string;
    const back = payload[FIELD_BACK] as string;
    const source = payload[FIELD_SOURCE] as string;
    const response = await HandlePostRequest(deckID, front, back, source);

    if (!response.ok) {
      throw new Error("Failed to create card");
    }

    toast.success("Card created successfully");
  } catch (error) {
    toast.error("Error submitting form");
  }
}

async function HandlePostRequest(
  deckID: string | undefined,
  front: string,
  back: string,
  source: string
) {
  return fetch(`http://localhost:8080/deck/${deckID}/card`, {
    method: "POST",
    body: JSON.stringify({
      content: {
        fields: ["front", "back"],
        values: [front, back],
      },
      source: source,
    }),
  });
}

export default CardFormPage;
