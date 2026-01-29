#include "utilspannels.h"

namespace pannelsBase
{
	void initPannel(
		QWidget* parent,
		QTreeView** listOfFiles,
		QLineEdit** curLocationEdit,
		QPushButton** backButton,
		QPushButton** forwardButton,
		QHeaderView** header,
		QHBoxLayout** pathSetting,
		QVBoxLayout** mainContainer)
	{
		*listOfFiles = new QTreeView(parent);
		(*listOfFiles)->setObjectName("listOfFiles");
		(*listOfFiles)->setSelectionBehavior(QAbstractItemView::SelectRows);
		(*listOfFiles)->setSelectionMode(QAbstractItemView::ExtendedSelection);
		(*listOfFiles)->setAlternatingRowColors(true);
		(*listOfFiles)->setUniformRowHeights(true);
		(*listOfFiles)->setFocusPolicy(Qt::StrongFocus);
		(*listOfFiles)->setItemsExpandable(false);
		(*listOfFiles)->setSortingEnabled(true);

		parent->setFocusProxy(*listOfFiles);

		*header = (*listOfFiles)->header();
		(*header)->setSectionsClickable(true);
		(*header)->setDefaultAlignment(Qt::AlignLeft);

		*backButton = new QPushButton("←", parent);
		(*backButton)->setFixedSize(30, 30);
		*forwardButton = new QPushButton("→", parent);
		(*forwardButton)->setFixedSize(30, 30);
		*curLocationEdit = new QLineEdit(parent);
		(*curLocationEdit)->setReadOnly(true);
		(*backButton)->setEnabled(false);
		(*forwardButton)->setEnabled(false);

		*pathSetting = new QHBoxLayout;
		(*pathSetting)->addWidget(*backButton);
		(*pathSetting)->addWidget(*forwardButton);
		(*pathSetting)->addWidget(*curLocationEdit);

		*mainContainer = new QVBoxLayout(parent);
		(*mainContainer)->addLayout(*pathSetting);
		(*mainContainer)->addWidget(*listOfFiles);
	}

	void sortByHeader(QAbstractItemModel* model, QTreeView* listOfFiles, int& lastHeaderFilterInd, SortingFilterWay& sortingWay, int headerInd)
	{
		if (!model || !listOfFiles)
			return;

		if (model->rowCount(QModelIndex()) == 0)
			return;

		if (lastHeaderFilterInd != headerInd)
		{
			lastHeaderFilterInd = headerInd;
			sortingWay = SortingFilterWay::Asc;
		}
		else
		{
			sortingWay = static_cast< SortingFilterWay >((static_cast< int >(sortingWay) + 1) % 3);
		}

		if (sortingWay == SortingFilterWay::None)
		{
			lastHeaderFilterInd = NOTHING_SORTED;
			model->sort(0, Qt::AscendingOrder);
		}
		else
		{
			Qt::SortOrder order = (sortingWay == SortingFilterWay::Asc) ? Qt::AscendingOrder : Qt::DescendingOrder;
			model->sort(headerInd, order);
		}
	}

}	 // namespace pannelsBase
