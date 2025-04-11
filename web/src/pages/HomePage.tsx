import { useEffect, useState } from "react";
import DeckList from "../components/DeckList";
import { DeckItem } from "../components/DeckList";
import NavBar from "../components/NavBar";
import { Link, useNavigate } from "react-router-dom";
import { toast } from "sonner";
import Spinner from "../components/Spinner";

function HomePage() {
  const [data, setData] = useState<DeckItem[] | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchData(setLoading, setError, setData);
  }, []);

  var body;

  if (loading) {
    body = <Spinner />;
  }

  if (error) {
    body = <></>;
  }

  if (data) {
    const itemList: DeckItem[] =
      data?.map((item) => ({
        id: item["id"] || "",
        name: item["name"] || "",
        description: item["description"] || "",
      })) || [];

    body = <DeckList items={itemList} />;
  }

  const rightComponents = [
    <Link to="/deck/form/-1" className="btn btn-light">
      Create New Deck
    </Link>,
  ];

  return (
    <>
      <NavBar rightComponents={rightComponents} />
      {body}
    </>
  );
}

async function fetchData(
  setLoading: (loading: boolean) => void,
  setError: (error: string | null) => void,
  setData: (data: DeckItem[] | null) => void
) {
  try {
    const response = await fetch("http://localhost:8080/deck");

    if (!response.ok) {
      throw new Error(`Http error! Status: ${response.status}`);
    }

    const result = await response.json();
    setData(result);
  } catch (error) {
    setError(`Error fetching data`);
    toast.error(`Error fetching data`);
  } finally {
    setLoading(false);
  }
}

export default HomePage;
