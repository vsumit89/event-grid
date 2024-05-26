import { getDateString } from '@/commons/dateTime';
import { getEvents } from '@/services/event';

import { useState, useEffect, useCallback } from 'react';

export const useEvents = (start, end) => {
  const [events, setEvents] = useState([]);

  const [loading, setLoading] = useState(true);

  const [error, setError] = useState('');

  const abortController = new AbortController();

  const fetchEvents = async () => {
    try {
        let startDate = getDateString(start);
        let endDate = getDateString(end);

        const data = await getEvents(startDate, endDate);
        
        setEvents(data);
    } catch (error) {
        setError(error.message);
    } finally {
        setLoading(false);
    }
  };

  useEffect(() => {
    fetchEvents();
    return () => {
      abortController.abort();
    };
  }, [start, end]);

  return { events, loading, error };
};

export default useEvents;