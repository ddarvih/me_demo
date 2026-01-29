#include "demo_fs_parser.h"

#include <QDebug>
#include <QTextCodec>

#define READ_UNTIL32 0x0FFFFFF8
#define READ_UNTIL16 0xFFF8

demo_fs_Parser::demo_fs_Parser()
{
	;
}

void demo_fs_Parser::getdemo_fs_Table()
{
	quint32 skipBeforeTable = reservedSectors * bytesInSector;
	file.seek(skipBeforeTable);

	quint32 tableInfoAmount = demo_fs_TableSize * bytesInSector;
	demo_fs_Table = file.read(tableInfoAmount);
}

bool demo_fs_Parser::unpackF(const QString& path)
{
	if (file.isOpen())
	{
		file.close();
		demo_fs_Table.clear();
	}

	file.setFileName(path);
	if (!file.open(QIODevice::ReadOnly))
		return false;

	QByteArray initInfo = file.read(512);
	if (initInfo.size() < 512)
	{
		file.close();
		return false;
	}
	if (initInfo.mid(510, 2) != QByteArray::fromHex("55AA"))
		return false;
	bytesInSector = static_cast< quint8 >(initInfo[11]) | (static_cast< quint8 >(initInfo[12]) << 8);
	if (bytesInSector % 512 != 0 || bytesInSector > 4096)
		return false;
	sectorsInCluster = static_cast< quint8 >(initInfo[13]);
	reservedSectors = static_cast< quint8 >(initInfo[14]) | (static_cast< quint8 >(initInfo[15]) << 8);
	demo_fs_TablesAmount = static_cast< quint8 >(initInfo[16]);
	quint16 demo_fs_16sz = static_cast< quint8 >(initInfo[22]) | (static_cast< quint8 >(initInfo[23]) << 8);
	quint32 countSectors;
	quint16 countSectors16 = static_cast< quint8 >(initInfo[19]) | (static_cast< quint8 >(initInfo[20]) << 8);
	countSectors = (countSectors16 != 0) ? countSectors16 : *reinterpret_cast< const quint32* >(initInfo.constData() + 32);
	rootEntries = static_cast< quint8 >(initInfo[17]) | (static_cast< quint8 >(initInfo[18]) << 8);
	quint32 countRootSectors = ((rootEntries * 32) + (bytesInSector - 1)) / bytesInSector;
	quint32 dataSectors = countSectors - (reservedSectors + (demo_fs_TablesAmount * demo_fs_TableSize) + countRootSectors);
	quint32 totalClusters = dataSectors / sectorsInCluster;
	isdemo_fs_32 = totalClusters >= 65525;

	quint32 valTemp;
	memcpy(&valTemp, initInfo.constData() + 36, sizeof(valTemp));
	demo_fs_TableSize = (isdemo_fs_32) ? valTemp : demo_fs_16sz;

	rootCluster = isdemo_fs_32 ? *reinterpret_cast< const quint32* >(initInfo.constData() + 44) : 0;
	quint32 rootDirSectors = ((rootEntries * 32) + (bytesInSector - 1)) / bytesInSector;
	startDataRegion = reservedSectors + (demo_fs_TablesAmount * demo_fs_TableSize);
	if (!isdemo_fs_32)
		startDataRegion += rootDirSectors;

	getdemo_fs_Table();
	return true;
}

quint64 demo_fs_Parser::findCurByCluster(quint32 cluster)
{
	if (cluster < 2)
		return 0;

	if (!isdemo_fs_32 && cluster == rootCluster)
	{
		return (reservedSectors + demo_fs_TablesAmount * demo_fs_TableSize) * bytesInSector;
	}
	return (startDataRegion + (cluster - 2) * sectorsInCluster) * bytesInSector;
}

QVector< quint32 > demo_fs_Parser::getClustersChain(quint32 start)
{
	QVector< quint32 > clusters;

	if (!isdemo_fs_32 && start == rootCluster)
	{
		clusters.append(0);
		return clusters;
	}

	quint32 cur = start;
	quint32 bytesPerOne = (isdemo_fs_32) ? 4 : 2;
	quint32 curSkip;
	while (true)
	{
		if (cur < 2)
			break;
		if (cur >= (isdemo_fs_32 ? READ_UNTIL32 : READ_UNTIL16))
			break;
		clusters.append(cur);

		curSkip = bytesPerOne * cur;
		if (curSkip + bytesPerOne > demo_fs_Table.size())
			break;
		cur = (isdemo_fs_32) ? *reinterpret_cast< const quint32* >(demo_fs_Table.constData() + curSkip) & 0x0FFFFFFF
						: *reinterpret_cast< const quint16* >(demo_fs_Table.constData() + curSkip);
	}
	return clusters;
}

QDateTime demo_fs_Parser::interpretTime(quint16 date, quint16 t)
{
	if (!date && !t)
		return QDateTime();

	int year = 1980 + ((date >> 9) & 0x7F);
	int mon = (date >> 5) & 0x0F;
	int day = date & 0x1F;
	QDate resDate = QDate(year, mon, day);

	int hour = (t >> 11) & 0x1F;
	int min = (t >> 5) & 0x3F;
	int sec = (t & 0x1F) * 2;
	QTime resTime = QTime(hour, min, sec);

	return QDateTime(resDate, resTime);
}

FileDescr demo_fs_Parser::parseEntry(const QByteArray& qba)
{
	FileDescr fd;

	if (qba.size() < 32)
		return FileDescr();

	QTextCodec* nameAdaptor = QTextCodec::codecForName("IBM 866");
	if (!nameAdaptor)
		return FileDescr();
	QByteArray namePart = qba.mid(0, 8);
	QByteArray extPart = qba.mid(8, 3);
	fd.name = nameAdaptor->toUnicode(namePart).trimmed() + nameAdaptor->toUnicode(extPart).trimmed();

	fd.isActiveItem = static_cast< quint8 >(qba[0]) != 0xE5;

	fd.attributes = static_cast< quint8 >(qba[11]);
	fd.type = (fd.attributes & 0x10) ? "Dir" : "File";

	if (isdemo_fs_32)
	{
		fd.clusterStart = ((quint16)(qba[20]) | ((quint16)(qba[21]) << 8)) << 16;
		fd.clusterStart |= (quint8)(qba[26]) | ((quint8)(qba[27]) << 8);
	}
	else
	{
		fd.clusterStart = (quint8)(qba[26]) | ((quint8)(qba[27]) << 8);
	}

	fd.size = *reinterpret_cast< const quint32* >(qba.constData() + 28);

	quint16 createT = *reinterpret_cast< const quint16* >(qba.constData() + 14);
	quint16 createD = *reinterpret_cast< const quint16* >(qba.constData() + 16);
	fd.creationTime = interpretTime(createD, createT);

	return fd;
}

QString demo_fs_Parser::getLName(const QVector< QByteArray >& entries)
{
	QString lName;
	for (const QByteArray& cur : entries)
	{
		for (const int pos : { 1, 3, 5, 7, 9, 14, 16, 18, 20, 22, 24, 28, 30 })
		{
			quint16 charCode = static_cast< quint8 >(cur[pos]) | (static_cast< quint8 >(cur[pos + 1]) << 8);
			if (charCode == 0x0000)
				return lName.trimmed();
			if (charCode == 0xFFFF)
				continue;
			lName += QChar(charCode);
		}
	}
	return lName.trimmed();
}

QVector< FileDescr > demo_fs_Parser::parseDir(quint32 start)
{
	QVector< FileDescr > dir;

	QVector< QByteArray > lfnEntrs;

	QVector< quint32 > dirClusters = getClustersChain(start);

	for (quint32 curClust : dirClusters)
	{
		quint64 moveToClust = findCurByCluster(curClust);
		if (!file.seek(moveToClust))
			continue;
		QByteArray clustInfo = file.read(bytesInSector * sectorsInCluster);

		for (int i = 0; i < clustInfo.size(); i += 32)
		{
			QByteArray entry = clustInfo.mid(i, 32);
			if (entry.isEmpty() || static_cast< quint8 >(entry[0]) == 0x00)
				break;

			if (static_cast< quint8 >(entry[11]) == 0x0F)
			{
				lfnEntrs.prepend(entry);
				continue;
			}

			FileDescr fd = parseEntry(entry);
			if (!lfnEntrs.isEmpty())
			{
				fd.name = getLName(lfnEntrs);
				lfnEntrs.clear();
			}
			dir.append(fd);
		}
	}
	return dir;
}

bool demo_fs_Parser::isOpen() const
{
	return file.isOpen();
}

void demo_fs_Parser::fillNeededDirInfo(QVector< FileDescr >& dir, const QByteArray& dirData, QVector< QByteArray >& lfn)
{
	for (int i = 0; i < dirData.size(); i += 32)
	{
		QByteArray entry = dirData.mid(i, 32);
		if (entry.isEmpty() || static_cast< quint8 >(entry[0]) == 0x00)
			break;

		if (isLfn(entry))
		{
			lfn.prepend(entry);
			continue;
		}

		FileDescr fd = parseEntry(entry);
		if (!lfn.isEmpty())
		{
			fd.name = getLName(lfn);
			lfn.clear();
		}

		if (fd.name == ".")
			continue;
		dir.append(fd);
	}
}

bool demo_fs_Parser::isLfn(const QByteArray& entry)
{
	return static_cast< quint8 >(entry[11]) == 0x0F;
}

QVector< FileDescr > demo_fs_Parser::dirList(quint32 cluster)
{
	QVector< FileDescr > wanted;
	QByteArray wData;
	if (!isdemo_fs_32 && cluster == 0)
	{
		quint16 goToRoot = bytesInSector * (reservedSectors + demo_fs_TablesAmount * demo_fs_TableSize);
		file.seek(goToRoot);
		quint32 haveToRead = bytesInSector * ((rootEntries * 32 + bytesInSector - 1) / bytesInSector);
		wData = file.read(haveToRead);
	}
	else
	{
		for (quint32 curClust : getClustersChain(cluster))
		{
			file.seek(findCurByCluster(curClust));
			wData.append(file.read(bytesInSector * sectorsInCluster));
		}
	}
	QVector< QByteArray > lfnInfo;
	fillNeededDirInfo(wanted, wData, lfnInfo);
	return wanted;
}

quint32 demo_fs_Parser::getParentCluster(quint32 cluster)
{
	if (cluster == 0 || (isdemo_fs_32 && cluster == rootCluster))
		return 0;

	for (const FileDescr& fd : dirList(cluster))
	{
		if (fd.type == "Dir" && fd.name == ".." && fd.isActiveItem)
		{
			return fd.clusterStart;
		}
	}
	return 0;
}

QByteArray demo_fs_Parser::readFileContent(const FileDescr& fd)
{
	if (fd.type == "Dir" || !fd.isActiveItem)
		return QByteArray();

	QByteArray readResult;
	quint64 resBytes = 0;
	QVector< quint32 > fileClusters = getClustersChain(fd.clusterStart);
	quint32 clusterBytesAmount = sectorsInCluster * bytesInSector;

	for (quint32 curClust : fileClusters)
	{
		quint64 moveToCurClust = findCurByCluster(curClust);
		if (!file.seek(moveToCurClust))
			return QByteArray();

		QByteArray clusterContent = file.read(clusterBytesAmount);

		quint64 bytesCount = qMin< quint32 >(clusterBytesAmount, fd.size - resBytes);
		readResult.append(clusterContent.left(bytesCount));
		resBytes += bytesCount;

		if (resBytes >= fd.size)
			break;
	}
	return readResult;
}

void demo_fs_Parser::close()
{
	if (file.isOpen())
		file.close();
	demo_fs_Table.clear();
	sectorsInCluster = 0;
	reservedSectors = 0;
	demo_fs_TablesAmount = 0;
	demo_fs_TableSize = 0;
	rootCluster = 0;
	isdemo_fs_32 = false;
	startDataRegion = 0;
}

demo_fs_Parser::~demo_fs_Parser()
{
	if (file.isOpen())
	{
		file.close();
	}
}

bool demo_fs_Parser::isdemo_fs_32Type() const
{
	return isdemo_fs_32;
}

quint32 demo_fs_Parser::getRootCluster() const
{
	return rootCluster;
}
