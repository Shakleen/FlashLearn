export interface DeckItem {
  id: string | number;
  name: string;
  description: string;
}

interface DeckListProps {
  items: DeckItem[];
}

function DeckList(props: DeckListProps) {
  const handleItemClick = (item: any) => {
    console.log("Item clicked:", item);
  };

  return (
    <>
      <ul className="list-group">
        {props.items.map((item) => (
          <li
            key={item.id}
            className="list-group-item list-group-item-action d-flex justify-content-between align-items-start"
            onClick={() => handleItemClick(item)}
            style={{ cursor: "pointer" }}
          >
            <div className="ms-2 me-auto">
              <div className="fw-bold">{item.name}</div>
              {item.description}
            </div>
          </li>
        ))}
      </ul>
    </>
  );
}

export default DeckList;
