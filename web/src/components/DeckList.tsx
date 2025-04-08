export interface DeckItem {
  id: string;
  name: string;
  description: string;
}

interface DeckListProps {
  items: DeckItem[];
}

function DeckList(props: DeckListProps) {
  return (
    <>
      <ul className="list-group">
        {props.items.map((item) => (
          <li className="list-group-item d-flex justify-content-between align-items-start">
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
