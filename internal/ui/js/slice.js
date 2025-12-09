document.addEventListener('DOMContentLoaded', function() {
  const slices = document.querySelectorAll('.slice');

  slices.forEach(function(slice) {
    const addBtn = slice.querySelector('.slice__add-btn');
    if (!addBtn) return;

    addBtn.addEventListener('click', function(e) {
      e.preventDefault();

      const items = slice.querySelectorAll('.slice__item');
      const template = items[items.length - 1];
      if (!template) return;

      const clone = template.cloneNode(true);
      const input = clone.querySelector('input');
      if (input) {
        const currentIndex = parseInt(input.name.split('.')[1]);
        const newIndex = currentIndex + 1;
        input.name = input.name.replace('.' + currentIndex, '.' + newIndex);
        input.value = '';
      }

      const removeBtn = clone.querySelector('.slice__remove-btn');
      if (removeBtn) {
        const currentIndex = parseInt(removeBtn.value.split(':')[2]);
        const newIndex = currentIndex + 1;
        removeBtn.value = removeBtn.value.replace(':' + currentIndex, ':' + newIndex);
      }

      template.parentNode.insertBefore(clone, addBtn);
    });

    slice.addEventListener('click', function(e) {
      if (e.target.classList.contains('slice__remove-btn')) {
        e.preventDefault();
        const item = e.target.closest('.slice__item');
        if (item) {
          item.remove();
        }
      }
    });
  });
});
