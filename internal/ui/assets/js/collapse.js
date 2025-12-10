function toggleCollapse(header) {
  const container = header.parentElement;
  container.classList.toggle('collapsed');
}

function toggleSliceFromLabel(label) {
  const field = label.parentElement;
  const slice = field.querySelector('.slice');
  if (slice) {
    slice.classList.toggle('collapsed');
  }
}
