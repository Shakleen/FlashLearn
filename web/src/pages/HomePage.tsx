import { useEffect, useState } from "react";
import DeckList from "../components/DeckList";
import { DeckItem } from "../components/DeckList";
import NavBar from "../components/NavBar";
import { useNavigate } from "react-router-dom";

function HomePage() {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(``);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch("http://localhost:8080/deck");

        if (!response.ok) {
          throw new Error(`Http error! Status: ${response.status}`);
        }

        const result = await response.json();
        setData(result);
      } catch (error) {
        setError(`Error fetching data`);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>{error}</div>;
  }

  const itemList: DeckItem[] = data.map((item) => ({
    id: item["id"] || "",
    name: item["name"] || "",
    description: item["description"] || "",
  }));

  return (
    <>
      <NavBar />
      <center>
        <button
          type="button"
          className="btn btn-primary m-3"
          onClick={() => {
            navigate("/deck/form/-1");
          }}
        >
          Create New Deck
        </button>
      </center>
      <DeckList items={itemList} />
    </>
  );
}

export default HomePage;
