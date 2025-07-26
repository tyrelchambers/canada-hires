export const getOpenTickets = (data: any) => {
  return data.filter((ticket: any) => ticket.status === "open");
};

export const getClosedTickets = (data: any) => {
  return data.filter((ticket: any) => ticket.status === "closed");
};

export const getPendingTickets = (data: any) => {
  return data.filter((ticket: any) => ticket.status === "pending");
};

export const getInProgressTickets = (data: any) => {
  return data.filter((ticket: any) => ticket.status === "in progress");
};
