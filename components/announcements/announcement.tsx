// TODO: better announcement type
import { announcementType } from '../../lib/types/types';

type announcementProps = {
  announcement: announcementType;
  editAnnouncement: (announcement: announcementType) => any;
  deleteAnnouncement: (id: string) => any;
};

const Announcement: React.FunctionComponent<announcementProps> = (props) => {
  const announcement = props.announcement;
  const editAnnouncement = props.editAnnouncement;
  const deleteAnnouncement = props.deleteAnnouncement;

  return (
    <div className="p-8 bg-blue-200 rounded-2xl shadow-lg">
      {announcement ? (
        <>
          <h1 className="text-2xl text-light">
            {announcement.title} ({announcement.id})
          </h1>
          <p className="text-md my-4">{announcement.description}</p>
          <a href={announcement.link} className="text-md my-4">Link</a>
          <div className="flex">
            <button
              className="mr-2 px-4 font-light rounded-lg bg-blue-300 hover:bg-blue-500"
              onClick={() => {
                editAnnouncement(announcement);
              }}
            >
              Edit
            </button>
            <button
              className="p-2 font-light rounded-lg bg-red-300 hover:bg-red-500"
              onClick={() => {
                deleteAnnouncement(announcement.id);
              }}
            >
              Delete
            </button>
          </div>
        </>
      ) : (
        <p>No information available.</p>
      )}
    </div>
  );
};

export default Announcement;
