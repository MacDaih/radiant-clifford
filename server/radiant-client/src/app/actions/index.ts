export const IS_LOADING = 'IS_LOADING';
export const DATA_RESPONSE = 'DATA_RESPONSE';

export function isLoading(isLoading = false) {
  return {
    type: IS_LOADING,
    isLoading: isLoading,
  };
}

function fetchDataResponse(json: JSON) {
  return {
    type: DATA_RESPONSE,
    json: json,
  };
}

export function fetchHello() {
  return (dispatch: any) => {
    dispatch(isLoading(true));
    fetch('http://localhost:8080/')
      .then(response => {
        return response;
      })
      .then(response => response.json())
      .then(json => {
        dispatch(isLoading(false));
        dispatch(fetchDataResponse(json));
      });
  };
}