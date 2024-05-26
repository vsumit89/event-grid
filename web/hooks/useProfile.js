import { getProfile } from '@/services/user';
import { useState, useEffect, useCallback } from 'react';

const useProfile = () => {
  const [profile, setProfile] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const abortController = new AbortController();

  const fetchProfile = useCallback(async () => {
    try {
        const profile = await getProfile();
        setProfile(profile);
    } catch (error) {
        setError(error.message);
    } finally {
        setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchProfile();
    return () => {
      abortController.abort();
    };
  }, [fetchProfile]);

  return { profile, loading, error };
};

export default useProfile;