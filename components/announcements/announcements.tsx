import React from 'react';
import { useState } from 'react';
import Announcement from './announcement';
import Modal from './modal';

const Announcements: React.FunctionComponent = () => {
  const [input, setInput] = useState('');
  const [data, setData] = useState([]);
  const [open, setOpen] = useState(null);
  const [method, setMethod] = useState('');
  const [message, setMessage] = useState('No results found.');

  const search = async () => {
    setMessage(`Searching database for ${input}...`);
    setData(await fetch(`/api/announcements/search?name=${input}`).then((res) => res.json()));
    setMessage('No results found.');
  };

  const editAnnouncement = (announcement) => {
    setMethod('PUT');
    setOpen(announcement);
  };

  const addModal = (announcement) => {
    setMethod('POST');
    setOpen({
      description: '',
      id: '',
      link: '',
      title: ''
    });
  };

  const deleteAnnouncement = async (id) => {
    await fetch(`/api/announcements/${id}`, {
      method: 'DELETE',
    }).then((res) => console.log(res.json()));
  };

  const close = async (announcement, update = false) => {
    if (update) {
      await fetch(`/api/announcements/${announcement.id}`, {
        method: method,
        body: JSON.stringify(announcement),
      })
        .then((res) => res.json())
        .then((msg) => console.log(msg))
        .catch((err) => console.error(err));
    }
    setOpen(null);
    setMethod('');
  };

  return (
    <div>
      {open ? <Modal info={open} close={close} /> : <></>}
      <div className="m-8">
        <div className="flex mb-8">
          <input
            value={input}
            className="ring-blue-200 mr-4 py-2 px-4 bg-white rounded-lg placeholder-gray-400 text-gray-900 appearance-none inline-block shadow-md focus:outline-none ring-2 focus:ring-blue-600"
            placeholder="search term"
            onInput={(e) => {
              const value = e.currentTarget.value;
              setInput(value);
            }}
            onKeyDown={(e) => {
              if (e.key === 'Enter') search();
            }}
          ></input>
          <button
            className="btn1"
            onClick={search}
          >
            Search
          </button>
          <button
            className="mx-4 btn1"
            onClick={addModal}
          >
            Add Announcement
          </button>
        </div>
        <div className="mt-4 grid grid-cols-1 gap-8">
          {data && data.length != 0 ? (
            data.map((announcement) => {
              return (
                <Announcement
                  key={announcement.id}
                  announcement={announcement}
                  editAnnouncement={editAnnouncement}
                  deleteAnnouncement={deleteAnnouncement}
                />
              );
            })
          ) : (
            <p className="text-center">{message}</p>
          )}
        </div>
      </div>
    </div>
  );
};

export default Announcements;
