import {useEffect} from 'react';

const useOutsideClick = (ref, action, params) => {
  useEffect(() => {
    function handle(event) {
      if (!ref.current?.contains(event.target)) {
        action(params);
      }
    }

    document.addEventListener('mousedown', handle);
    return () => { document.removeEventListener('mousedown', handle); };
  }, [ref, action, params]);
};

export default useOutsideClick;