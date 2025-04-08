import DeckList from "./components/DeckList";
import { DeckItem } from "./components/DeckList";

function App() {
  const itemList: DeckItem[] = [
    { id: "0", name: "Deck #0", description: "Description #0" },
    { id: "1", name: "Deck #1", description: "Description #1" },
    { id: "2", name: "Deck #2", description: "Description #2" },
    { id: "3", name: "Deck #3", description: "Description #3" },
    { id: "4", name: "Deck #4", description: "Description #4" },
  ];

  return (
    <div>
      <DeckList items={itemList} />
    </div>
  );
}

export default App;
