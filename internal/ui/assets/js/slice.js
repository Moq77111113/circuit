document.addEventListener("DOMContentLoaded", function () {
  const slices = document.querySelectorAll(".slice");

  slices.forEach(function (slice) {
    const addbutton = slice.querySelector(".slice__add-button");
    if (!addbutton) return;

    addbutton.addEventListener("click", function (e) {
      e.preventDefault();

      const items = slice.querySelectorAll(".slice__item");
      const template = items[items.length - 1];
      if (!template) return;

      const clone = template.cloneNode(true);
      const input = clone.querySelector("input");
      if (input) {
        const currentIndex = parseInt(input.name.split(".")[1]);
        const newIndex = currentIndex + 1;
        input.name = input.name.replace("." + currentIndex, "." + newIndex);
        input.value = "";
      }

      const removebutton = clone.querySelector(".slice__remove-button");
      if (removebutton) {
        const currentIndex = parseInt(removebutton.value.split(":")[2]);
        const newIndex = currentIndex + 1;
        removebutton.value = removebutton.value.replace(
          ":" + currentIndex,
          ":" + newIndex
        );
      }

      template.parentNode.insertBefore(clone, addbutton);
    });

    slice.addEventListener("click", function (e) {
      if (e.target.classList.contains("slice__remove-button")) {
        e.preventDefault();
        const item = e.target.closest(".slice__item");
        if (item) {
          item.remove();
        }
      }
    });
  });
});
